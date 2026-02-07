class Retrotxt < Formula
  desc "Convert and display legacy text files and ANSI art on modern terminals"
  homepage "https://github.com/bengarrett/retrotxtgo"
  url "https://github.com/bengarrett/retrotxtgo/archive/refs/tags/v1.1.1.tar.gz"
  sha256 "f845137f5061f7e92414f42a283eca8fa23a6c01457c2675917e5a73ac72b2de"
  version "1.1.1"
  license "LGPL-3.0-only"

  @commit = "cbf2357109ae17fbb0d4dcff59a4dee6bc486dd8"
  @build_date = "2026-02-07T22:54:10+11:00"

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