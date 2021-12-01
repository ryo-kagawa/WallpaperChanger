// https://docs.microsoft.com/ja-jp/windows/win32/winprog/windows-data-types

package windows

const (
	FALSE BOOL = 0
	TRUE  BOOL = 1
)

type BOOL int32
type DWORD uint32
type PVOID uintptr
type HANDLE PVOID
type HDC HANDLE
type HMONITOR HANDLE
type HRESULT LONG
type HWND HANDLE
type LONG int32
type LONG_PTR int64
type LPARAM LONG_PTR
type LPCWSTR *WCHAR
type LPVOID uintptr
type LPWSTR *WCHAR
type UINT uint32
type WCHAR uint16
