class Semrel < Formula
  desc "Semantic Versioning and Conventional Commits release utility"
  homepage "https://github.com/greatliontech/semrel"
  version "VERSION"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/greatliontech/semrel/releases/download/v#{version}/semrel_#{version}_darwin_arm64.tar.gz"
      sha256 "CHECKSUM_DARWIN_ARM64"
    else
      url "https://github.com/greatliontech/semrel/releases/download/v#{version}/semrel_#{version}_darwin_amd64.tar.gz"
      sha256 "CHECKSUM_DARWIN_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/greatliontech/semrel/releases/download/v#{version}/semrel_#{version}_linux_arm64.tar.gz"
      sha256 "CHECKSUM_LINUX_ARM64"
    else
      url "https://github.com/greatliontech/semrel/releases/download/v#{version}/semrel_#{version}_linux_amd64.tar.gz"
      sha256 "CHECKSUM_LINUX_AMD64"
    end
  end

  def install
    bin.install "semrel"

    generate_completions_from_executable(bin/"semrel", "completion")
  end

  test do
    system "#{bin}/semrel", "--help"
  end
end
