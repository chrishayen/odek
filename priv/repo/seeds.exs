alias Valkyrie.Accounts.User
alias Valkyrie.Repo

admin_email = "chris@shotgun.dev"
admin_password = "changeme"

admin_attrs = %{
  email: admin_email,
  hashed_password: Bcrypt.hash_pwd_salt(admin_password),
  confirmed_at: DateTime.utc_now(:second),
  is_admin: true,
  must_change_password: true
}

case Repo.get_by(User, email: admin_email) do
  nil ->
    %User{}
    |> Ecto.Changeset.change(admin_attrs)
    |> Repo.insert!()

  user ->
    user
    |> Ecto.Changeset.change(admin_attrs)
    |> Repo.update!()
end
