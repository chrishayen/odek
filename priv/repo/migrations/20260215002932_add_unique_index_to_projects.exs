defmodule Valkyrie.Repo.Migrations.AddUniqueIndexToProjects do
  use Ecto.Migration

  def change do
    create unique_index(:projects, [:organization_id, "lower(name)"],
             where: "deleted_at is null",
             name: :projects_org_lower_name_uniq
           )
  end
end
