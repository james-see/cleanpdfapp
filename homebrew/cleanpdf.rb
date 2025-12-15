# Homebrew Cask formula for Clean PDF
# Copy this file to your homebrew-tap repository: Casks/cleanpdf.rb
#
# Installation: brew install james-see/tap/cleanpdf
# Or: brew tap james-see/tap && brew install cleanpdf

cask "cleanpdf" do
  version "2.0.0"
  sha256 :no_check # Update with actual SHA256 after first release

  url "https://github.com/james-see/cleanpdfapp/releases/download/v#{version}/CleanPDF-macos-v#{version}.zip"
  name "Clean PDF"
  desc "View and wipe metadata from PDF files"
  homepage "https://james-see.github.io/cleanpdfapp/"

  livecheck do
    url :url
    strategy :github_latest
  end

  app "CleanPDF.app"

  zap trash: [
    "~/Library/Application Support/us.jamescampbell.cleanpdf",
    "~/Library/Caches/us.jamescampbell.cleanpdf",
    "~/Library/Preferences/us.jamescampbell.cleanpdf.plist",
  ]
end

