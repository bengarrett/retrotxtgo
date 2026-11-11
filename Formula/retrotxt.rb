class Retrotxt < Formula
  desc "Convert and display legacy text files and ANSI art on modern terminals"
  homepage "https://github.com/bengarrett/retrotxtgo"
  url "https://github.com/bengarrett/retrotxtgo/archive/refs/tags/v1.2.5.tar.gz"
  sha256 "f0189058c0cb192577e572a7893888061ede25e843af84fd50406df86dce7edb"
  version "1.2.5"
  license "LGPL-3.0-only"

  @commit = "d3f2f28a80412fa7a51b83ad64da1d4c83097ece"
  @build_date = "2026-11-12T08:33:31+11:00"

  livecheck do
    url :stable
    strategy :github_latest
  end

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version} -X main.commit=#{self.class.instance_variable_get('@commit')} -X main.date=#{self.class.instance_variable_get('@build_date')}")
  end

  test do
    assert_match "retrotxt", shell_output("#{bin}/retrotxt --version")
  end
end
