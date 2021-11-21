package user32

var fWinIni = struct {
	SPIF_UPDATEINIFILE uintptr
	SPIF_SENDCHANGE    uintptr
}{
	SPIF_UPDATEINIFILE: 0x01,
	SPIF_SENDCHANGE:    0x02,
}
