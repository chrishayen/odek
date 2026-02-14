defmodule Valkyrie.Prompts do
  @moduledoc """
  Skills/rules/profiles and system prompt composition.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.Accounts.Scope

  alias Valkyrie.Prompts.{
    ProjectRuleProfile,
    ProjectSkillProfile,
    Rule,
    RuleProfile,
    RuleProfileItem,
    Skill,
    SkillProfile,
    SkillProfileItem,
    SystemPrompt,
    SystemPromptVersion
  }

  alias Valkyrie.Repo
  alias Valkyrie.Workspace.Project

  @default_prompt_keys ~w(chat_planning worker)

  def create_skill(%Scope{} = scope, attrs) do
    attrs
    |> normalize_attrs()
    |> Map.merge(%{"organization_id" => scope.organization_id, "owner_user_id" => scope.user.id})
    |> then(&(%Skill{} |> Skill.changeset(&1) |> Repo.insert()))
  end

  def create_rule(%Scope{} = scope, attrs) do
    attrs
    |> normalize_attrs()
    |> Map.merge(%{"organization_id" => scope.organization_id, "owner_user_id" => scope.user.id})
    |> then(&(%Rule{} |> Rule.changeset(&1) |> Repo.insert()))
  end

  def create_skill_profile(%Scope{} = scope, attrs) do
    attrs
    |> normalize_attrs()
    |> Map.merge(%{"organization_id" => scope.organization_id, "owner_user_id" => scope.user.id})
    |> then(&(%SkillProfile{} |> SkillProfile.changeset(&1) |> Repo.insert()))
  end

  def create_rule_profile(%Scope{} = scope, attrs) do
    attrs
    |> normalize_attrs()
    |> Map.merge(%{"organization_id" => scope.organization_id, "owner_user_id" => scope.user.id})
    |> then(&(%RuleProfile{} |> RuleProfile.changeset(&1) |> Repo.insert()))
  end

  def add_skill_profile_item(%Scope{} = scope, profile_id, skill_id) do
    with {:ok, _profile} <- fetch_skill_profile(scope.organization_id, profile_id),
         {:ok, _skill} <- fetch_skill(scope.organization_id, skill_id) do
      %SkillProfileItem{}
      |> SkillProfileItem.changeset(%{profile_id: profile_id, skill_id: skill_id})
      |> Repo.insert(on_conflict: :nothing)
      |> normalize_insert_conflict()
    end
  end

  def add_rule_profile_item(%Scope{} = scope, profile_id, rule_id) do
    with {:ok, _profile} <- fetch_rule_profile(scope.organization_id, profile_id),
         {:ok, _rule} <- fetch_rule(scope.organization_id, rule_id) do
      %RuleProfileItem{}
      |> RuleProfileItem.changeset(%{profile_id: profile_id, rule_id: rule_id})
      |> Repo.insert(on_conflict: :nothing)
      |> normalize_insert_conflict()
    end
  end

  def assign_skill_profile(%Scope{} = scope, project_id, profile_id) do
    with {:ok, _project} <- fetch_project(scope.organization_id, project_id),
         {:ok, _profile} <- fetch_skill_profile(scope.organization_id, profile_id) do
      %ProjectSkillProfile{}
      |> ProjectSkillProfile.changeset(%{project_id: project_id, profile_id: profile_id})
      |> Repo.insert(on_conflict: :nothing)
      |> normalize_insert_conflict()
    end
  end

  def assign_rule_profile(%Scope{} = scope, project_id, profile_id) do
    with {:ok, _project} <- fetch_project(scope.organization_id, project_id),
         {:ok, _profile} <- fetch_rule_profile(scope.organization_id, profile_id) do
      %ProjectRuleProfile{}
      |> ProjectRuleProfile.changeset(%{project_id: project_id, profile_id: profile_id})
      |> Repo.insert(on_conflict: :nothing)
      |> normalize_insert_conflict()
    end
  end

  def list_profile_kinds do
    ["skill_profiles", "rule_profiles"]
  end

  def list_prompt_keys(%Scope{} = scope) do
    db_keys =
      Repo.all(
        from p in SystemPrompt,
          where: p.organization_id == ^scope.organization_id,
          order_by: [asc: p.prompt_key],
          select: p.prompt_key
      )

    (db_keys ++ @default_prompt_keys)
    |> Enum.uniq()
    |> Enum.sort()
  end

  def add_prompt_version(%Scope{} = scope, prompt_key, version, body) do
    Repo.transact(fn ->
      prompt =
        Repo.one(
          from p in SystemPrompt,
            where: p.organization_id == ^scope.organization_id and p.prompt_key == ^prompt_key,
            lock: "FOR UPDATE"
        )

      prompt =
        case prompt do
          nil ->
            {:ok, created} =
              %SystemPrompt{}
              |> SystemPrompt.changeset(%{
                organization_id: scope.organization_id,
                prompt_key: prompt_key,
                name: prompt_key,
                active_version: version
              })
              |> Repo.insert()

            created

          existing ->
            existing
        end

      case %SystemPromptVersion{}
           |> SystemPromptVersion.changeset(%{
             system_prompt_id: prompt.id,
             version: version,
             body: body,
             created_by_user_id: scope.user.id
           })
           |> Repo.insert(
             on_conflict: [set: [body: body, created_by_user_id: scope.user.id]],
             conflict_target: [:system_prompt_id, :version]
           ) do
        {:ok, _prompt_version} -> :ok
        {:error, changeset} -> Repo.rollback(changeset)
      end

      if is_nil(prompt.active_version) do
        {:ok, _updated} =
          prompt
          |> SystemPrompt.changeset(%{active_version: version})
          |> Repo.update()
      end

      %{prompt_key: prompt_key, version: version}
    end)
    |> normalize_tx_result()
  end

  def activate_prompt_version(%Scope{} = scope, prompt_key, version) do
    with %SystemPrompt{} = prompt <-
           Repo.one(
             from p in SystemPrompt,
               where: p.organization_id == ^scope.organization_id and p.prompt_key == ^prompt_key
           ),
         %SystemPromptVersion{} <-
           Repo.one(
             from pv in SystemPromptVersion,
               where: pv.system_prompt_id == ^prompt.id and pv.version == ^version
           ),
         {:ok, updated_prompt} <-
           prompt |> SystemPrompt.changeset(%{active_version: version}) |> Repo.update() do
      {:ok, updated_prompt}
    else
      nil -> {:error, :not_found}
      error -> error
    end
  end

  def compose_prompt(%Scope{} = scope, project_id, prompt_key, user_prompt \\ nil) do
    layers = [
      active_system_prompt_body(scope.organization_id, prompt_key),
      global_skill_bodies(scope.organization_id, scope.user.id),
      project_skill_profile_bodies(project_id),
      global_rule_bodies(scope.organization_id, scope.user.id),
      project_rule_profile_bodies(project_id),
      optional_user_prompt(user_prompt)
    ]

    layers
    |> List.flatten()
    |> Enum.filter(&(is_binary(&1) and String.trim(&1) != ""))
    |> Enum.join("\n\n---\n\n")
  end

  defp active_system_prompt_body(organization_id, prompt_key) do
    prompt =
      Repo.one(
        from p in SystemPrompt,
          where: p.organization_id == ^organization_id and p.prompt_key == ^prompt_key
      )

    case prompt do
      nil ->
        ""

      %SystemPrompt{active_version: nil} ->
        ""

      %SystemPrompt{} = prompt ->
        Repo.one(
          from pv in SystemPromptVersion,
            where: pv.system_prompt_id == ^prompt.id and pv.version == ^prompt.active_version,
            select: pv.body
        ) || ""
    end
  end

  defp global_skill_bodies(organization_id, owner_user_id) do
    Repo.all(
      from s in Skill,
        where:
          s.organization_id == ^organization_id and s.owner_user_id == ^owner_user_id and
            is_nil(s.deleted_at),
        order_by: [asc: s.inserted_at, asc: s.id],
        select: s.body
    )
  end

  defp project_skill_profile_bodies(project_id) do
    Repo.all(
      from psp in ProjectSkillProfile,
        join: spi in SkillProfileItem,
        on: spi.profile_id == psp.profile_id,
        join: s in Skill,
        on: s.id == spi.skill_id,
        where: psp.project_id == ^project_id and is_nil(s.deleted_at),
        order_by: [asc: psp.inserted_at, asc: spi.inserted_at, asc: s.inserted_at, asc: s.id],
        select: s.body
    )
  end

  defp global_rule_bodies(organization_id, owner_user_id) do
    Repo.all(
      from r in Rule,
        where:
          r.organization_id == ^organization_id and r.owner_user_id == ^owner_user_id and
            is_nil(r.deleted_at),
        order_by: [asc: r.inserted_at, asc: r.id],
        select: r.body
    )
  end

  defp project_rule_profile_bodies(project_id) do
    Repo.all(
      from prp in ProjectRuleProfile,
        join: rpi in RuleProfileItem,
        on: rpi.profile_id == prp.profile_id,
        join: r in Rule,
        on: r.id == rpi.rule_id,
        where: prp.project_id == ^project_id and is_nil(r.deleted_at),
        order_by: [asc: prp.inserted_at, asc: rpi.inserted_at, asc: r.inserted_at, asc: r.id],
        select: r.body
    )
  end

  defp optional_user_prompt(value) when is_binary(value) and value != "", do: [value]
  defp optional_user_prompt(_value), do: []

  defp fetch_project(org_id, project_id) do
    case Repo.one(
           from p in Project,
             where: p.id == ^project_id and p.organization_id == ^org_id and is_nil(p.deleted_at)
         ) do
      nil -> {:error, :not_found}
      project -> {:ok, project}
    end
  end

  defp fetch_skill_profile(org_id, profile_id) do
    case Repo.one(
           from p in SkillProfile,
             where: p.id == ^profile_id and p.organization_id == ^org_id and is_nil(p.deleted_at)
         ) do
      nil -> {:error, :not_found}
      profile -> {:ok, profile}
    end
  end

  defp fetch_rule_profile(org_id, profile_id) do
    case Repo.one(
           from p in RuleProfile,
             where: p.id == ^profile_id and p.organization_id == ^org_id and is_nil(p.deleted_at)
         ) do
      nil -> {:error, :not_found}
      profile -> {:ok, profile}
    end
  end

  defp fetch_skill(org_id, skill_id) do
    case Repo.one(
           from s in Skill,
             where: s.id == ^skill_id and s.organization_id == ^org_id and is_nil(s.deleted_at)
         ) do
      nil -> {:error, :not_found}
      skill -> {:ok, skill}
    end
  end

  defp fetch_rule(org_id, rule_id) do
    case Repo.one(
           from r in Rule,
             where: r.id == ^rule_id and r.organization_id == ^org_id and is_nil(r.deleted_at)
         ) do
      nil -> {:error, :not_found}
      rule -> {:ok, rule}
    end
  end

  defp normalize_attrs(attrs) do
    Map.new(attrs, fn {k, v} -> {to_string(k), v} end)
  end

  defp normalize_tx_result({:ok, value}), do: {:ok, value}
  defp normalize_tx_result({:error, reason}), do: {:error, reason}

  defp normalize_insert_conflict({:ok, struct}), do: {:ok, struct}
  defp normalize_insert_conflict(other), do: other
end
