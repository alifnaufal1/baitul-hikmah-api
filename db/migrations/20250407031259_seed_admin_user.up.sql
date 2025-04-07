INSERT INTO users (username, password, role)
VALUES (
  'adminganteng',
  '$2a$08$tezVhSP8qe23txwAJkQ82OwtOHVbAN2Ok.FZ1mt98hPpCw8TpfflC',
  'admin'
)
ON CONFLICT (username) DO NOTHING;
