{pkgs, ...}:
pkgs.buildGoModule rec {
  pname = "dashdashdash";
  version = "unstable";
  src = builtins.path {
    name = "${pname}-source";
    path = ./.;
  };

  ldflags = ["-X main.version=${version}"];
  buildTarget = "./cmd";

  postInstall = ''
    mv $out/bin/cmd $out/bin/${pname}
  '';

  vendorHash = "sha256-KMDWQgShV8+amolOUy/NZpoF/6bHTuujrYosM1NcinQ=";
  doCheck = true;
}
