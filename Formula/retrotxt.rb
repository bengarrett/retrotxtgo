class Retrotxt < Formula
  desc "Convert and display legacy text files and ANSI art on modern terminals"
  homepage "https://github.com/bengarrett/retrotxtgo"
  url "https://github.com/bengarrett/retrotxtgo/archive/refs/tags/v1.2.2.tar.gz"
  sha256 "604e49734ade8f54fc9f05bd8d623464fdac6c5d3de3ee27c7637f6353340e62"
  version "1.2.2"
  license "LGPL-3.0-only"

  @commit = "3c2466a99a4e9366d0d1bdf0bc2cbde2937703fd"
  @build_date = "2026-08-07T23:08:38+10:00"

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
