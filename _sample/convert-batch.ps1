$Dir = "C:/Users/kaede/Downloads/35167"
$Bin = "C:/Users/kaede/src/github.com/ledyba/pdf-sharpener/pdf-sharpener.exe"

$Files = Get-ChildItem $Dir -Name -Recurse -Include *.orig.pdf
foreach ($File in $Files){
  $Dst = $File.Replace(".orig.pdf", ".pdf")
  echo $File
  $Input = Join-Path $Dir $File
  $Output = Join-Path $Dir $Dst
  & $Bin -i $Input -o $Output
}
