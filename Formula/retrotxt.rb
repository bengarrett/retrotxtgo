class Retrotxt < Formula
  desc "Convert and display legacy text files and ANSI art on modern terminals"
  homepage "https://github.com/bengarrett/retrotxtgo"
  url "https://github.com/bengarrett/retrotxtgo/archive/refs/tags/v1.2.1.tar.gz"
  sha256 "c1d2e32f5b868974672ae8eeb1c675f5214e28c1e323ffcb880c9e3f4e3039c8"
  version "1.2.1"
  license "LGPL-3.0-only"

  @commit = "90da32819cbe9c31f93d04e03c6f799c109751b9"
  @build_date = "2026-05-01T10:53:49+10:00"

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
