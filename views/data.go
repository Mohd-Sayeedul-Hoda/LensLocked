package views

type Data struct{
  Alert *Alert
  Yield interface{}
}

type Alert struct{
  Level string
  Message string
}

const(
  AlertLvlError = "danger"
  AlertLvlWarning = "warning"
  AlertLvlInfo = "info"
  AlertLvlSuccess = "success"
)
