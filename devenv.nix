{ pkgs, lib, config, inputs, ... }:

let
  postgres_address = "127.0.0.1";
  postgres_port = 5432;
  postgres_database = "shorter";
  postgres_sslmode = "disable";
in
{
  cachix.enable = false;

  # https://devenv.sh/basics/
  env.GREET = "Shorter";
  env.SHORTER_POSTGRES_ADDRESS = postgres_address;
  env.SHORTER_POSTGRES_PORT = postgres_port;
  env.SHORTER_POSTGRES_DATABASE = postgres_database;
  env.SHORTER_POSTGRES_SSLMODE = postgres_sslmode;

  packages = with pkgs; [ templ ];

  scripts.hello.exec = ''
    echo -e "
     🚀 Welcome to the devenv for \e[32m$GREET\e[0m 🔥
    make sure to run \e[34mdevenv up\e[0m in a seperate terminal to start postgresql
    "
  '';

  enterShell = ''
    echo ""
    git -v
    go version
    templ version | sed -e "s/^/Templ version: /;"
    hello
  '';

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep "2.42.0"
  '';

  # https://devenv.sh/services/
  # services.postgres.enable = true;
  services = {
    postgres = {
      enable = true;
      listen_addresses = "localhost";
      initialDatabases = [{ name = postgres_database; }];
    };
    adminer.enable = true;
  };

  # https://devenv.sh/languages/
  # languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}
