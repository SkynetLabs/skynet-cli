class SkynetCli < Formula
  desc "Lightweight CLI to interact with Skynet"
  homepage "https://siasky.net"
  url "https://github.com/SkynetLabs/skynet-cli.git"
  version "2.1.0"
  sha256 "0cad3711fa115a5c0ff746cb42806362dcf5ca7c4fd944cde59fbbddb607c6ea"
  license "MIT"

  depends_on "go" => :build

  def install
    build_time = Utils.safe_popen_read("date").chomp
    git_revision = Utils.safe_popen_read("git", "rev-parse", "--short", "HEAD").chomp
    git_dirty = Utils.safe_popen_read("git", "diff-index", "--quiet", "HEAD", "--", "||", "echo", "'x-'").chomp
    ldflags = %W[
      -s -w
      -X github.com/SkynetLabs/skynet-cli/build.GitRevision=#{git_dirty}#{git_revision}
      -X "github.com/SkynetLabs/skynet-cli/build.BuildTime=#{build_time}"
    ].join(" ")
    system "go", "build", *std_go_args, "-tags", "netgo", "-ldflags", ldflags, "./cmd/skynet"
    mv bin/"skynet-cli", bin/"skynet"
  end

  test do
    str_version = shell_output("#{bin}/skynet version")
    assert_match "skynet #{version}", str_version

    str_help = shell_output("#{bin}/skynet help")
    str_default = shell_output("#{bin}/skynet")

    assert_match str_version, str_help
    assert_match str_default.lines[0], str_help
    assert_match "Usage:", str_help
    assert_match "Available Commands:", str_help

    err_test = shell_output("#{bin}/skynet upload foo")
    expected_message = "no such file or directory"
    assert_match expected_message, err_test
  end
end
