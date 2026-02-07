class Retrotxt < Formula
  desc "Convert and display legacy text files and ANSI art on modern terminals"
  homepage "https://github.com/bengarrett/retrotxtgo"
  url "https://github.com/bengarrett/retrotxtgo/archive/refs/tags/v1.2.0.tar.gz"
  sha256 "8aa0f0dac9d53fad1f8d36f37bf47a6af18ab0edc7874aa8e4ea2d9ca7328d32"
  version "1.2.0"
  license "LGPL-3.0-only"

  @commit = "f433b468d2a4854c3ed5a1db77f5405d5dc4934b"
  @build_date = "2026-02-08T00:15:16+11:00"

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