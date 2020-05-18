package tests

import "log"

// LogoASCII is the RetroTxt, CP-437 ASCII stored as base64.
const LogoASCII = ""

// LogoANSI is the RetroTxt, CP-437, 24-bit color ANSI stored as base64.
const LogoANSI = "G1swOzQwOzM3bQ0KG1szNW0bWzE7MTkxOzA7NXQg2xtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ03NwbWzE7MjU1OzE1NTsxNTd03BtbMTsyNTU7MTAyOzEwNnTc3BtbMDszNW0bWzE7MTkxOzA7NXTb2xtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ03NwbWzE7MjU1OzE1NTsxNTd03NzcG1sxOzI1NTsxMDI7MTA2dNzcG1sxOzI1NTsxNTU7MTU3dNwbWzE7MjU1OzEwMjsxMDZ03NwbWzA7MzVtG1sxOzE5MTswOzV029sbWzE7NDVtG1swOzE5MTswOzV0G1sxOzI1NTsxMDI7MTA2dNzcG1swOzM1bRtbMTsxOTE7MDs1dNuy29vb29vf39/f398gICAgG1sxbRtbMTsyMzg7MTA0OzI1M3Tf3BtbMDszNW0bWzE7MTkxOzA7NXQgIN/fG1sxbRtbMTsyMzg7MTA0OzI1M3TcG1swOzM1bRtbMTsxOTE7MDs1dN/f39/b29vbstsbWzE7NDVtG1swOzE5MTswOzV0G1sxOzI1NTsxMDI7MTA2dNzcG1sxOzI1NTsxNTU7MTU3dNwbWzE7MjU1OzEwMjsxMDZ03BtbMDszNW0bWzE7MTkxOzA7NXTb2xtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ03BtbMTsyNTU7MTU1OzE1N3Tc3BtbMTsyNTU7MTAyOzEwNnTc3BtbMDszNW0bWzE7MTkxOzA7NXTb2xtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ03NzcG1sxOzI1NTsxNTU7MTU3dNzc3BtbMTsyNTU7MTAyOzEwNnTc3BtbMDszNW0bWzE7MTkxOzA7NXTbDQobWzM2bRtbMTswOzIxNzsyMzR0IBtbMzVtG1sxOzE5MTswOzV029/f39/f39/f39/f39/f39/f39/f39/f3xtbMzdtG1sxNkMbWzE7MzVtG1sxOzI0NTsxNzY7MjU1dN/c3N8bWzA7MzVtG1sxOzE2ODsyOzE1NXQgICAbWzM2bRtbMTswOzIxNzsyMzR0IBtbMzdtG1s1QxtbMzVtG1sxOzE5MTswOzV039/f39/f39/f39/f39/f39/f39/f39+yDQobWzM2bRtbMTswOzIxNzsyMzR0IBtbMzVtG1sxOzE5MTswOzV029uyG1sxOzQ1bRtbMDsxOTE7MDs1dBtbMTsyNTU7MTAyOzEwNnTf3xtbNW0bWzA7MjU1OzEwMjsxMDZ0G1sxOzI1NTsxNTU7MTU3dN/fG1swOzE7NDU7MzVtG1swOzE5MTswOzV0G1sxOzI1NTsxMDI7MTA2dN/fG1swOzM1bRtbMTsxOTE7MDs1dNvbG1sxOzQ1bRtbMDsxOTE7MDs1dBtbMTsyNTU7MTAyOzEwNnTf3xtbNW0bWzA7MjU1OzEwMjsxMDZ0G1sxOzI1NTsxNTU7MTU3dN8bWzA7MTs0NTszNW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ0398bWzA7MzVtG1sxOzE5MTswOzV029vb29vf3xtbMTsxNjg7MjsxNTV0ICAbWzFtG1sxOzI0NTsxNzY7MjU1dNzf3xtbMDszNW0bWzE7MjMwOzA7MjI0dN8bWzQ1bRtbMDsxNjg7MjsxNTV0shtbNDBtG1sxOzE2ODsyOzE1NXQgG1s0NW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0shtbNDBt3xtbMW0bWzE7MjQ1OzE3NjsyNTV039/fG1swOzM1bRtbMTsyMzA7MDsyMjR03xtbMW0bWzE7MjQ1OzE3NjsyNTV03xtbMTsyMzg7MTA0OzI1M3Tf3xtbMDszNW0bWzE7MjMwOzA7MjI0dN8bWzFtG1sxOzI0NTsxNzY7MjU1dN8bWzU7NDU7MzZtG1swOzI0NTsxNzY7MjU1dBtbMTsxMjg7MjU1OzI1NXTfG1swOzE7MzVtG1sxOzI0NTsxNzY7MjU1dN/cG1sxOzIzODsxMDQ7MjUzdN8bWzA7MzVtG1sxOzIzMDswOzIyNHTfG1sxOzQ1bRtbMDsyMzA7MDsyMjR0G1sxOzIzODsxMDQ7MjUzdNsbWzA7MzVtG1sxOzE2ODsyOzE1NXQgG1sxOzIzMDswOzIyNHTb3xtbMW0bWzE7MjM4OzEwNDsyNTN039/cG1swOzM1bRtbMTsxNjg7MjsxNTV0ICAbWzE7MTkxOzA7NXTf39vbshtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ0398bWzA7MzVtG1sxOzE5MTswOzV029sbWzE7NDVtG1swOzE5MTswOzV0G1sxOzI1NTsxMDI7MTA2dN/f3xtbNW0bWzA7MjU1OzEwMjsxMDZ0G1sxOzI1NTsxNTU7MTU3dN/fG1swOzE7NDU7MzVtG1swOzE5MTswOzV0G1sxOzI1NTsxMDI7MTA2dN/f3xtbMDszNW0bWzE7MTkxOzA7NXTb27Lb2w0KG1szNm0bWzE7MDsyMTc7MjM0dCAbWzM1bRtbMTsxOTE7MDs1dNzcG1sxbRtbMTsyNTU7MTAyOzEwNnTc3NzcG1swOzM1bRtbMTsxOTE7MDs1dNwbWzFtG1sxOzI1NTsxMDI7MTA2dNzcG1swOzM1bRtbMTsxOTE7MDs1dNwbWzFtG1sxOzI1NTsxMDI7MTA2dNzcG1swOzM1bRtbMTsxOTE7MDs1dNzc3NwbWzM3bRtbNkMbWzE7MzVtG1sxOzIzODsxMDQ7MjUzdNwbWzE7MjQ1OzE3NjsyNTV03N8bWzA7MzJtG1sxOzA7MTQ3OzN03BtbNDJtG1swOzA7MTQ3OzN0G1sxOzA7MjA0OzB03BtbMW0bWzE7NDU7MjU1OzExM3TcG1swOzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLEbWzQwbRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLEbWzQwOzMybRtbMTswOzE0NzszdLIbWzE7NDJtG1swOzA7MjA0OzB0G1sxOzQ1OzI1NTsxMTN03BtbMDs0MjszMm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHTcG1s0MG0bWzE7MDsxNDc7M3TfIN8bWzQybRtbMDswOzE0NzszdBtbMTswOzIwNDswdN8bWzQwbRtbMTswOzE0NzszdN8bWzE7MzVtG1sxOzIzODsxMDQ7MjUzdNwbWzE7MjQ1OzE3NjsyNTV03xtbMDszMm0bWzE7MDsxNDc7M3QgINwbWzE7MzVtG1sxOzI0NTsxNzY7MjU1dN8bWzA7MzJtG1sxOzA7MTQ3OzN0shtbMTs0NTszNW0bWzA7MjMwOzA7MjI0dBtbMTsyMzg7MTA0OzI1M3SyG1swOzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dNsbWzQwOzMybRtbMTswOzE0NzszdLIbWzE7NDJtG1swOzA7MTQ3OzN0G1sxOzQ1OzI1NTsxMTN03BtbMDs0MjszMm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHTcG1s0MG0bWzE7MDsxNDc7M3TcG1szNW0bWzE7MjMwOzA7MjI0dN8bWzFtG1sxOzI0NTsxNzY7MjU1dNwbWzE7MjM4OzEwNDsyNTN03BtbMDszNW0bWzE7MTY4OzI7MTU1dCAgICAbWzM3bRtbNUMbWzM1bRtbMTsxOTE7MDs1dNzc3BtbMW0bWzE7MjU1OzEwMjsxMDZ03BtbMTsyNTU7MTU1OzE1N3TcG1sxOzI1NTsxMDI7MTA2dNwbWzA7MzVtG1sxOzE5MTswOzV03BtbMW0bWzE7MjU1OzEwMjsxMDZ03NzcG1swOzM1bRtbMTsxOTE7MDs1dNzcDQobWzM2bRtbMTswOzIxNzsyMzR0IBtbMzVtG1sxOzE5MTswOzV039/f39/f39/f398bWzFtG1sxOzIzODsxMDQ7MjUzdNsbWzA7MzVtG1sxOzE5MTswOzV0398bWzE7MTY4OzI7MTU1dCAbWzFtG1sxOzI0NTsxNzY7MjU1dCAbWzA7MzVtG1sxOzE2ODsyOzE1NXQgIBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dNsbWzQwbd8bWzFtG1sxOzI0NTsxNzY7MjU1dN8bWzA7MzVtG1sxOzIzMDswOzIyNHTfG1szMm0bWzE7MDsxNDc7M3Tc3N/fG1s0Mm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHSx3xtbNDA7MzVtG1sxOzE2ODsyOzE1NXQgG1s0NW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sBtbNDBtG1sxOzE2ODsyOzE1NXQgG1s0Nm0bWzA7MDsxNDU7MTQxdLIbWzQwOzMybRtbMTswOzE0NzszdLAbWzQybRtbMDswOzE0NzszdBtbMTswOzIwNDswdLDfG1s0MDszNW0bWzE7MTY4OzI7MTU1dLAgsRtbNDI7MzJtG1swOzA7MTQ3OzN0G1sxOzA7MjA0OzB027IbWzQwbdwbWzM1bRtbMTsxNjg7MjsxNTV0sCCwG1s0MjszMm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHSy3xtbNDBtG1sxOzA7MTQ3OzN0sBtbMTs0NTszNW0bWzA7MjMwOzA7MjI0dBtbMTsyMzg7MTA0OzI1M3SwG1swOzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLIbWzQwOzMybRtbMTswOzE0NzszdLEbWzQybRtbMDswOzE0NzszdBtbMTswOzIwNDswdN+yG1s0MG3fG1sxOzA7MTQ3OzN03xtbMTswOzIwNDswdNwbWzE7MDsxNDc7M3TcG1sxOzM1bRtbMTsyMzg7MTA0OzI1M3Tf3xtbMDszNW0bWzE7MjMwOzA7MjI0dN8bWzQ1bRtbMDsxNjg7MjsxNTV0shtbNDBtG1sxOzE2ODsyOzE1NXQgG1sxOzIzMDswOzIyNHQgICAbWzE7MTY4OzI7MTU1dCAgICAbWzE7MTkxOzA7NXTf39/f39/f398NChtbMzZtG1sxOzA7MjE3OzIzNHQgG1szNW0bWzE7MTkxOzA7NXTbshtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ03xtbMTsyNTU7MTU1OzE1N3TfG1sxOzI1NTsxMDI7MTA2dN8bWzA7MzVtG1sxOzE5MTswOzV02xtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ03xtbMDszNW0bWzE7MTkxOzA7NXTb3xtbMW0bWzE7MjM4OzEwNDsyNTN03xtbMDszNW0bWzE7MjMwOzA7MjI0dCAbWzFtG1sxOzI0NTsxNzY7MjU1dNsgG1swOzM1bRtbMTsyMzA7MDsyMjR03N/fG1s0NW0bWzA7MTY4OzI7MTU1dLIbWzQwbRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLIbWzQwOzMybRtbMTswOzE0NzszdNwbWzE7NDJtG1swOzA7MTQ3OzN0G1sxOzQ1OzI1NTsxMTN03BtbMDszMm0bWzE7MDsxNDc7M3TbG1szNW0bWzE7MTY4OzI7MTU1dCDc39sbWzMybRtbMTswOzE0NzszdNuyG1szNW0bWzE7MTY4OzI7MTU1dCAbWzQ1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSwG1s0MG0bWzE7MTY4OzI7MTU1dCAbWzQ2bRtbMDswOzE0NTsxNDF0sRtbNDBtIBtbMzJtG1sxOzA7MTQ3OzN0srIbWzM1bRtbMTsxNjg7MjsxNTV0siCyG1s0MjszMm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHSxsLIbWzQwOzM1bRtbMTsxNjg7MjsxNTV0siCyG1s0MjszMm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHSwsBtbNDA7MzVtG1sxOzE2ODsyOzE1NXQgG1s0NW0bWzA7MjMwOzA7MjI0dLAbWzQwbSAbWzQ1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSxG1s0MDszMm0bWzE7MDsxNDc7M3SwstsbWzM1bRtbMTsxNjg7MjsxNTV029zcIBtbMzJtG1sxOzA7MTQ3OzN03xtbNDJtG1swOzA7MTQ3OzN0G1sxOzA7MjA0OzB03xtbNDA7MzVtG1sxOzE2ODsyOzE1NXQgG1s0NW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sRtbNDBtG1sxOzE2ODsyOzE1NXQgG1sxOzIzMDswOzIyNHQgG1sxbRtbMTsyMzg7MTA0OzI1M3Tc3xtbMTsyNDU7MTc2OzI1NXTfG1sxOzIzODsxMDQ7MjUzdN/c3BtbMTsyNDU7MTc2OzI1NXTcG1swOzM1bRtbMTsxNjg7MjsxNTV0ICAbWzE7MTkxOzA7NXTf3xtbMTs0NW0bWzA7MTkxOzA7NXQbWzE7MjU1OzEwMjsxMDZ0398bWzA7MzVtG1sxOzE5MTswOzV029sNChtbMzZtG1sxOzA7MjE3OzIzNHQgG1szNW0bWzE7MTkxOzA7NXTc3Nzc3CAbWzE7MTY4OzI7MTU1dCAgG1sxbRtbMTsyNDU7MTc2OzI1NXTfG1s0NW0bWzA7MTY4OzI7MTU1dN/fG1s0MDszNm0bWzE7MTI4OzI1NTsyNTV02xtbMzVtG1sxOzI0NTsxNzY7MjU1dN/f3xtbMDszNW0bWzE7MTY4OzI7MTU1dCAbWzQ1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSxG1s0MG0bWzE7MTY4OzI7MTU1dCAbWzQ1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSxG1s0MDszMm0bWzE7MDsxNDc7M3SyG1s0Mm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHTf2xtbNDA7MzVtG1sxOzE2ODsyOzE1NXTbICDbG1szMm0bWzE7MDsxNDc7M3SysRtbMzVtG1sxOzE2ODsyOzE1NXQg2xtbMzZtG1sxOzA7MTQ1OzE0MXQgG1s0NjszNW0bWzA7MDsxNDU7MTQxdBtbMTsxNjg7MjsxNTV0sBtbNDA7MzZtG1sxOzA7MTQ1OzE0MXQgG1szMm0bWzE7MDsxNDc7M3SxsBtbNDY7MzVtG1swOzA7MTQ1OzE0MXQbWzE7MTY4OzI7MTU1dLIbWzQwbSDbG1szMm0bWzE7MDsxNDc7M3Sy2xtbNDJtG1swOzA7MTQ3OzN0G1sxOzA7MjA0OzB0sBtbNDA7MzVtG1sxOzE2ODsyOzE1NXTbINsbWzMybRtbMTswOzE0NzszdLKyG1szNW0bWzE7MTY4OzI7MTU1dCAbWzQ1bRtbMDsyMzA7MDsyMjR0sRtbNDBtIBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLAbWzQwbRtbMTsxNjg7MjsxNTV0IBtbMzJtG1sxOzA7MTQ3OzN0sbIbWzM1bRtbMTsxNjg7MjsxNTV0siAg2xtbMzJtG1sxOzA7MTQ3OzN027IbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLAbWzQwbRtbMTsxNjg7MjsxNTV0IBtbMTsyMzA7MDsyMjR02yAbWzE7MzJtG1sxOzQ1OzI1NTsxMTN03BtbMDs0MjszMm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHTcG1sxbRtbMTs0NTsyNTU7MTEzdNwbWzA7MzJtG1sxOzA7MTQ3OzN03NwbWzM1bRtbMTsxNjg7MjsxNTV0IBtbMTsyMzA7MDsyMjR03xtbMW0bWzE7MjM4OzEwNDsyNTN03xtbMDszNW0bWzE7MjMwOzA7MjI0dNwbWzE7MTY4OzI7MTU1dCAgG1sxOzE5MTswOzV03NzcDQobWzM2bRtbMTswOzIxNzsyMzR0IBtbMzVtG1sxOzE5MTswOzV0398bWzE7MTY4OzI7MTU1dCAgG1sxOzIzMDswOzIyNHTcG1sxbRtbMTsyMzg7MTA0OzI1M3TcG1swOzM1bRtbMTsyMzA7MDsyMjR03BtbMW0bWzE7MjM4OzEwNDsyNTN03xtbMDszNW0bWzE7MjMwOzA7MjI0dN8bWzMwbSAgG1sxOzM1bRtbMTsyNDU7MTc2OzI1NXTbG1sxOzIzODsxMDQ7MjUzdCDfG1swOzMybRtbMTswOzE0NzszdNsbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLAbWzQwbRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLAbWzQwOzMybRtbMTswOzE0NzszdLEbWzQybRtbMDswOzE0NzszdBtbMTswOzIwNDswdLGyG1s0NTszNW0bWzA7MjMwOzA7MjI0dBtbMTsxNjg7MjsxNTV0shtbNDBtICAbWzQ2bRtbMDswOzE0NTsxNDF0shtbNDA7MzJtG1sxOzA7MTQ3OzN0sbAbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDZtG1swOzA7MTQ1OzE0MXSyG1s0MDszNm0bWzE7MDsxNDU7MTQxdCAg39zcG1s0NjszNW0bWzA7MDsxNDU7MTQxdBtbMTsxNjg7MjsxNTV0sRtbNDA7MzZtG1sxOzA7MTQ1OzE0MXQgG1s0NjszNW0bWzA7MDsxNDU7MTQxdBtbMTsxNjg7MjsxNTV0shtbNDA7MzJtG1sxOzA7MTQ3OzN0sbKyG1s0NTszNm0bWzA7MTY4OzI7MTU1dBtbMTswOzE0NTsxNDF0sBtbNDA7MzVtG1sxOzE2ODsyOzE1NXQg2xtbMzJtG1sxOzA7MTQ3OzN0sbEbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzIzMDswOzIyNHSyG1s0MG0g2yAbWzMybRtbMTswOzE0NzszdLCxG1szNW0bWzE7MTY4OzI7MTU1dLIgILIbWzMybRtbMTswOzE0NzszdLKxG1szNW0bWzE7MTY4OzI7MTU1dCDbIBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLAbWzQwOzMybRtbMTswOzE0NzszdLEbWzQybRtbMDswOzE0NzszdBtbMTswOzIwNDswdN+yG1s0MG0bWzE7MDsxNDc7M3Tf3xtbNDJtG1swOzA7MTQ3OzN0G1sxOzA7MjA0OzB03xtbMW0bWzE7NDU7MjU1OzExM3TcG1swOzQyOzMybRtbMDswOzE0NzszdBtbMTswOzIwNDswdNwbWzQwbRtbMTswOzE0NzszdNwbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLIbWzQwbRtbMTsxNjg7MjsxNTV0IBtbMzZtG1sxOzA7MTQ1OzE0MXQgG1szNW0bWzE7MTkxOzA7NXTf3w0KG1szNm0bWzE7MDsyMTc7MjM0dCAbWzM1bRtbMTsxOTE7MDs1dLLb3xtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLIbWzQwOzMybRtbMTswOzE0NzszdNwbWzQybRtbMDswOzE0NzszdBtbMTswOzIwNDswdNwbWzFtG1sxOzQ1OzI1NTsxMTN03BtbMDszMm0bWzE7MDsxNDc7M3TfG1szNW0bWzE7MTY4OzI7MTU1dNzfG1sxbRtbMTsyNDU7MTc2OzI1NXSyG1swOzMybRtbMTswOzE0NzszdLKyshtbMzVtG1sxOzE2ODsyOzE1NXQg2yDbG1szMm0bWzE7MDsxNDc7M3Sw2xtbNDJtG1swOzA7MTQ3OzN0G1sxOzA7MjA0OzB0sRtbNDU7MzVtG1swOzIzMDswOzIyNHQbWzE7MTY4OzI7MTU1dLIbWzQwbSAgG1s0Nm0bWzA7MDsxNDU7MTQxdLEbWzQwOzMybRtbMTswOzE0NzszdLAbWzM2bRtbMTswOzE0NTsxNDF0IBtbMzVtG1sxOzE2ODsyOzE1NXQgG1s0Nm0bWzA7MDsxNDU7MTQxdLAbWzQwOzM2bRtbMTswOzE0NTsxNDF0IBtbMTswOzEwMzsxMDh03BtbMzdtG1s1QxtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLEbWzQwOzMybRtbMTswOzE0NzszdLCxsRtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLEbWzQwOzM1bRtbMTsxNjg7MjsxNTV0IBtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLAbWzQwOzMybRtbMTswOzE0NzszdLCwG1szNW0bWzE7MTY4OzI7MTU1dCDbIBtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLAbWzQwOzM1bRtbMTsxNjg7MjsxNTV0ICAbWzMybRtbMTswOzE0NzszdLAbWzM1bRtbMTsxNjg7MjsxNTV02yDc2xtbMzJtG1sxOzA7MTQ3OzN0sbAbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLAbWzQwOzM1bRtbMTsxNjg7MjsxNTV0INsbWzMybRtbMTswOzE0NzszdLCyG1s0Mm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHSwG1s0MDszNW0bWzE7MTY4OzI7MTU1dNvf3NwbWzQyOzMybRtbMDswOzE0NzszdBtbMTswOzIwNDswdLKxG1s0MG0bWzE7MDsxNDc7M3SwG1s0NTszNW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sRtbNDBtG1sxOzE2ODsyOzE1NXQgG1sxOzE5MTswOzV0IN/fDQobWzM2bRtbMTswOzIxNzsyMzR0IBtbMzVtG1sxOzE5MTswOzV03CAbWzE7MTY4OzI7MTU1dCAgG1s0NW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sRtbNDA7MzJtG1sxOzA7MTQ3OzN0shtbNDJtG1swOzA7MTQ3OzN0G1sxOzA7MjA0OzB0sN8bWzQwOzM1bRtbMTsxNjg7MjsxNTV0shtbMzZtG1sxOzA7MTQ1OzE0MXQgG1szNW0bWzE7MTY4OzI7MTU1dCCyG1szMm0bWzE7MDsxNDc7M3SxsbEbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDZtG1swOzA7MTQ1OzE0MXSyG1s0MG0g2yAbWzMybRtbMTswOzE0NzszdLLbG1szNW0bWzE7MTY4OzI7MTU1dNsgG1szNm0bWzE7MDsxNDU7MTQxdCAbWzQ2OzM1bRtbMDswOzE0NTsxNDF0G1sxOzE2ODsyOzE1NXSwG1s0MDszNm0bWzE7MDsxNDU7MTQxdNzc398gIBtbMTswOzEwMzsxMDh0stvc3BtbMTswOzE0NTsxNDF0IBtbNDVtG1swOzE2ODsyOzE1NXSyG1s0MDszNW0bWzE7MTY4OzI7MTU1dCAbWzMybRtbMTswOzE0NzszdLCwG1s0NTszNm0bWzA7MTY4OzI7MTU1dBtbMTswOzE0NTsxNDF0shtbNDBtIBtbNDVtG1swOzE2ODsyOzE1NXSxG1s0MDszNW0bWzE7MTY4OzI7MTU1dCAgIBtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLAbWzQwOzM1bRtbMTsxNjg7MjsxNTV0IBtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLEbWzQwOzM1bRtbMTsxNjg7MjsxNTV0ICAbWzMybRtbMTswOzE0NzszdLAbWzM1bRtbMTsxNjg7MjsxNTV039/fIBtbMzJtG1sxOzA7MTQ3OzN0sBtbMzVtG1sxOzE2ODsyOzE1NXQgG1s0NTszNm0bWzA7MTY4OzI7MTU1dBtbMTswOzE0NTsxNDF0sRtbNDA7MzVtG1sxOzE2ODsyOzE1NXQgIBtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLAbWzQwOzM1bRtbMTsxNjg7MjsxNTV0IBtbMzJtG1sxOzA7MTQ3OzN0sbIbWzQ2OzM1bRtbMDswOzE0NTsxNDF0G1sxOzE2ODsyOzE1NXSyG1s0MG0gIBtbNDZtG1swOzA7MTQ1OzE0MXSyG1s0MDszMm0bWzE7MDsxNDc7M3SyG1s0Mm0bWzA7MDsxNDc7M3QbWzE7MDsyMDQ7MHSwG1s0MDszNW0bWzE7MTY4OzI7MTU1dCAbWzQ1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSwG1s0MG0bWzE7MTY4OzI7MTU1dCAbWzM2bRtbMTswOzEwMzsxMDh0sBtbMzVtG1sxOzE5MTswOzV0IN8NChtbMzZtG1sxOzA7MjE3OzIzNHQgIBtbMzVtG1sxOzE2ODsyOzE1NXQgG1szNm0bWzE7MDsxMDM7MTA4dLAbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLAbWzQwOzMybRtbMTswOzE0NzszdLGyshtbMzVtG1sxOzE2ODsyOzE1NXSyICDbG1szMm0bWzE7MDsxNDc7M3SwsLAbWzM2bRtbMTswOzE0NTsxNDF0IBtbNDY7MzVtG1swOzA7MTQ1OzE0MXQbWzE7MTY4OzI7MTU1dLEbWzQwOzM2bRtbMTswOzE0NTsxNDF0IBtbNDY7MzVtG1swOzA7MTQ1OzE0MXQbWzE7MTY4OzI7MTU1dLIbWzQwbSAbWzMybRtbMTswOzE0NzszdLGyG1szNW0bWzE7MTY4OzI7MTU1dNsbWzM3bRtbMTBDG1szNm0bWzE7MDsxMDM7MTA4dLLbG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzA7MTAzOzEwOHQbWzE7MDsyMTc7MjM0dLAbWzQwbRtbMTswOzE0NTsxNDF0ICAgG1s0Nm0bWzA7MDsxMDM7MTA4dBtbMTswOzIxNzsyMzR0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgG1s0NW0bWzA7MTY4OzI7MTU1dLIbWzQwbSAgG1szNW0bWzE7MTY4OzI7MTU1dCAbWzQ1OzM2bRtbMDsxNjg7MjsxNTV0G1sxOzA7MTQ1OzE0MXSxG1s0MDszNW0bWzE7MTY4OzI7MTU1dCAbWzQ1OzM2bRtbMDsxNjg7MjsxNTV0G1sxOzA7MTQ1OzE0MXSyG1s0MG0gG1sxOzA7MTAzOzEwOHQgG1szNW0bWzE7MTY4OzI7MTU1dCAgIBtbMzZtG1sxOzA7MTQ1OzE0MXQgINzfIBtbMTswOzEwMzsxMDh03BtbMTswOzE0NTsxNDF0IBtbNDVtG1swOzE2ODsyOzE1NXSxG1s0MDszNW0bWzE7MTY4OzI7MTU1dCAbWzMybRtbMTswOzE0NzszdLCxG1s0NjszNW0bWzA7MDsxNDU7MTQxdBtbMTsxNjg7MjsxNTV0sRtbNDBtICAbWzQ2bRtbMDswOzE0NTsxNDF0sRtbNDA7MzJtG1sxOzA7MTQ3OzN0sbIbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDZtG1swOzA7MTQ1OzE0MXSyG1swbQ0KG1szNm0bWzE7MDsyMTc7MjM0dCAbWzE7MDsxMDM7MTA4dLAbWzE7MDsyMTc7MjM0dCAbWzM1bRtbMTsxNjg7MjsxNTV0ICDbG1szMm0bWzE7MDsxNDc7M3SwsbEbWzM1bRtbMTsxNjg7MjsxNTV02yDc3xtbMzJtG1sxOzA7MTQ3OzN0sBtbMzZtG1sxOzA7MTQ1OzE0MXQgIBtbNDY7MzVtG1swOzA7MTQ1OzE0MXQbWzE7MTY4OzI7MTU1dLAbWzQwOzM2bRtbMTswOzE0NTsxNDF0ICAbWzQ1bRtbMDsxNjg7MjsxNTV0sBtbNDA7MzVtG1sxOzE2ODsyOzE1NXQgG1szMm0bWzE7MDsxNDc7M3SwsRtbMzVtG1sxOzE2ODsyOzE1NXTf39/fG1s0Nm0bWzA7MDsxNDU7MTQxdLIbWzQwOzM2bRtbMTswOzE0NTsxNDF0ICAbWzE7MDsxMDM7MTA4dLCwsN/f3BtbMTswOzE0NTsxNDF0ICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzMybRtbMTswOzE0NzszdLCwG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sBtbNDBtG1sxOzA7MTQ1OzE0MXQg2yAgIBtbNDVtG1swOzE2ODsyOzE1NXSyG1s0MG0g2yAgIBtbMTswOzEwMzsxMDh039wbWzQ2bRtbMDswOzE0NTsxNDF03BtbNDBtG1sxOzA7MTQ1OzE0MXTfIBtbMTswOzEwMzsxMDh03NywG1sxOzA7MTQ1OzE0MXQgG1s0NW0bWzA7MTY4OzI7MTU1dLIbWzQwbSAgG1szMm0bWzE7MDsxNDc7M3SwG1s0NjszNW0bWzA7MDsxNDU7MTQxdBtbMTsxNjg7MjsxNTV0sBtbNDBtICAbWzQ2bRtbMDswOzE0NTsxNDF0sBtbNDA7MzJtG1sxOzA7MTQ3OzN0sLEbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDZtG1swOzA7MTQ1OzE0MXSxG1s0MG0gG1szNm0bWzE7MDsxMDM7MTA4dLANChtbMTswOzIxNzsyMzR0IBtbMTswOzEwMzsxMDh0sbCwG1szNW0bWzE7MTY4OzI7MTU1dCAbWzQ1OzM2bRtbMDsxNjg7MjsxNTV0G1sxOzA7MTQ1OzE0MXSwG1s0MDszNW0bWzE7MTY4OzI7MTU1dCAgG1szMm0bWzE7MDsxNDc7M3SwG1szNW0bWzE7MTY4OzI7MTU1dN/fG1szNm0bWzE7MDsyMTc7MjM0dCAgG1sxOzA7MTQ1OzE0MXQg3N8gICAbWzQ1bRtbMDsxNjg7MjsxNTV0sRtbNDA7MzVtG1sxOzE2ODsyOzE1NXQgG1szNm0bWzE7MDsxNDU7MTQxdCAbWzMybRtbMTswOzE0NzszdLAbWzM2bRtbMTswOzE0NTsxNDF0ICAg3N8gIBtbMTswOzEwMzsxMDh02xtbMzdtG1s1QxtbMzZtG1sxOzA7MTAzOzEwOHTcG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0shtbMTs0MDszNW0bWzE7MjM4OzEwNDsyNTN0sBtbMDszNm0bWzE7MDsxNDU7MTQxdCAgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sRtbNDBtG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sBtbMTs0MDszNW0bWzE7MjM4OzEwNDsyNTN0sLAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQg2yAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSwG1s0MDszN20bWzVDG1szNm0bWzE7MDsxMDM7MTA4dN/f3BtbMTswOzE0NTsxNDF0IBtbMTswOzEwMzsxMDh0398gG1sxOzA7MTQ1OzE0MXTbICAg2yAbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDZtG1swOzA7MTQ1OzE0MXSwG1s0MG0gG1szMm0bWzE7MDsxNDc7M3SwG1szNW0bWzE7MTY4OzI7MTU1dCAbWzQ2bRtbMDswOzE0NTsxNDF0sBtbNDBtICAgG1szNm0bWzE7MDsxMDM7MTA4dLANChtbMTswOzIxNzsyMzR0IBtbMTswOzEwMzsxMDh0srEbWzE7MDsxNDU7MTQxdCAbWzM1bRtbMTsxNjg7MjsxNTV0IBtbNDU7MzZtG1swOzE2ODsyOzE1NXQbWzE7MDsxNDU7MTQxdLEbWzQwOzM1bRtbMTsxNjg7MjsxNTV0IBtbMzZtG1sxOzA7MjE3OzIzNHQgICAbWzM1bRtbMTsxNjg7MjsxNTV0INwbWzM2bRtbMTswOzE0NTsxNDF0398gG1sxOzA7MTAzOzEwOHSw3N8bWzE7MDsxNDU7MTQxdCAbWzQ1bRtbMDsxNjg7MjsxNTV0shtbNDBtICAgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTswOzEwMzsxMDh0shtbNDBt3xtbMTswOzE0NTsxNDF03yAbWzE7MDsxMDM7MTA4dNzc37CwsLCwsBtbMTswOzE0NTsxNDF0ICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHTbG1sxOzQwOzM1bRtbMTsyMzg7MTA0OzI1M3SwG1swOzM2bRtbMTswOzE0NTsxNDF0ICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSxG1sxOzQwOzM1bRtbMTsyMzg7MTA0OzI1M3SwsBtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSxG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MG0bWzE7MDsxNDU7MTQxdCAgG1sxOzA7MTAzOzEwOHQgICAgG1sxOzA7MTQ1OzE0MXQgIBtbMTswOzEwMzsxMDh039/cG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgIBtbMzVtG1sxOzIzMDswOzIyNHSwG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sRtbNDBtG1sxOzA7MTQ1OzE0MXQgINsgICDbICAbWzE7MDsyMTc7MjM0dCAbWzE7MDsxMDM7MTA4dLENChtbMTswOzIxNzsyMzR0IBtbMTswOzEwMzsxMDh027KwG1szNW0bWzE7MTY4OzI7MTU1dCAbWzQ1OzM2bRtbMDsxNjg7MjsxNTV0G1sxOzA7MTQ1OzE0MXSyG1s0MDszNW0bWzE7MTY4OzI7MTU1dLCwIBtbMzZtG1sxOzA7MjE3OzIzNHQgIBtbMTswOzE0NTsxNDF0IBtbMzVtG1sxOzE2ODsyOzE1NXTf3NwgIBtbMzZtG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLAbWzQwOzM1bRtbMTsxNjg7MjsxNTV0sLAbWzM2bRtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MDsxMDM7MTA4dLEbWzQwbRtbMTswOzE0NTsxNDF0ICAbWzE7MDsxMDM7MTA4dNuy39/fG1szN20bWzdDG1s0NjszNm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0sBtbMTs0MDszNW0bWzE7MjM4OzEwNDsyNTN0sbCwG1swOzQ2OzM2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHTbG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MG3c3N8bWzE7MDsxNDU7MTQxdCAgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB02xtbNDA7MzVtG1sxOzIzMDswOzIyNHSwG1szNm0bWzE7MDsxNDU7MTQxdCAgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sRtbNDBtG1sxOzA7MTQ1OzE0MXTf3NwgICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzA7MTAzOzEwOHSxG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MDszNW0bWzE7MjMwOzA7MjI0dLCwsRtbNDY7MzZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdNsbWzQwbRtbMTswOzE0NTsxNDF0ICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzM1bRtbMTsyMzA7MDsyMjR0sBtbMzZtG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHSwG1sxOzA7MjE3OzIzNHQgG1sxOzA7MTAzOzEwOHSyDQobWzE7MDsyMTc7MjM0dCAgG1sxOzA7MTAzOzEwOHTbshtbMTswOzE0NTsxNDF0IBtbNDVtG1swOzE2ODsyOzE1NXTbG1s0MDszNW0bWzE7MTY4OzI7MTU1dLGxsCAgG1szNm0bWzE7MDsyMTc7MjM0dCAgG1sxOzA7MTQ1OzE0MXQgIBtbMzVtG1sxOzE2ODsyOzE1NXTf39wbWzM2bRtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLEbWzQwOzM1bRtbMTsxNjg7MjsxNTV0sbCwG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTswOzEwMzsxMDh0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgICDc3N/fG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHSwsLDbG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0sRtbMTs0MDszNW0bWzE7MjM4OzEwNDsyNTN0srGxG1swOzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dLCwsNzfG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0sBtbNDA7MzVtG1sxOzIzMDswOzIyNHSwsBtbMzZtG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sRtbNDBtG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLAbWzQwbRtbMTswOzE0NTsxNDF0ICAgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTswOzEwMzsxMDh0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0sBtbNDA7MzVtG1sxOzIzMDswOzIyNHSxsbIbWzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0shtbNDBtG1sxOzA7MTQ1OzE0MXQgG1szNW0bWzE7MjMwOzA7MjI0dLCwG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sRtbNDBtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHSx27INChtbMTswOzIxNzsyMzR0IBtbMTswOzEwMzsxMDh029zbG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sBtbNDA7MzVtG1sxOzE2ODsyOzE1NXTbsrIbWzQ1OzM2bRtbMDsxNjg7MjsxNTV0G1sxOzA7MTQ1OzE0MXTbG1s0MG3f39wgG1szNW0bWzE7MTY4OzI7MTU1dLCxsRtbNDY7MzZtG1swOzA7MTQ1OzE0MXQbWzE7MDsxMDM7MTA4dLAbWzQwbRtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdNsbWzQwOzM1bRtbMTsxNjg7MjsxNTV0srGxG1szNm0bWzE7MDsxNDU7MTQxdNsgIBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MDsxMDM7MTA4dLAbWzQwbRtbMTswOzE0NTsxNDF0IBtbMzVtG1sxOzE2ODsyOzE1NXSwsLAbWzQ2OzM2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MG0bWzE7MDsxNDU7MTQxdCAgIBtbMTswOzEwMzsxMDh027IbWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSyG1sxOzQwOzM1bRtbMTsyMzg7MTA0OzI1M3TbsrIbWzA7NDY7MzZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLIbWzQwbRtbMTswOzE0NTsxNDF0ICAbWzE7MDsxMDM7MTA4dNvb3xtbMTswOzE0NTsxNDF0ICAbWzQ2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSyG1s0MDszNW0bWzE7MjMwOzA7MjI0dLGwG1szNm0bWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MG0bWzE7MDsxNDU7MTQxdCAgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0shtbNDA7MzVtG1sxOzIzMDswOzIyNHSwsBtbMzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzM1bRtbMTsyNTU7MTU1OzE1N3QgG1swOzM2bRtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLEbWzQwOzM1bRtbMTsyMzA7MDsyMjR0srLbG1s0NjszNm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0sRtbNDBtG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdNsbWzQwOzM1bRtbMTsyMzA7MDsyMjR0sLGxG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0shtbNDBtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHSy298NChtbMTswOzIxNzsyMzR0IBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MDsxMDM7MTA4dLIbWzQwbdzfG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sRtbNDU7MzVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLAbWzQwbRtbMTsxNjg7MjsxNTV027IbWzQ2OzM2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dLAbWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSwG1s0MDszNW0bWzE7MTY4OzI7MTU1dLLb2xtbNDY7MzZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLEbWzQwbRtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLAbWzQwOzM1bRtbMTsxNjg7MjsxNTV027IbWzQ1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSwG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLAbWzQwOzM1bRtbMTsxNjg7MjsxNTV0sLGxsBtbNDY7MzZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLAbWzQwbRtbMTswOzE0NTsxNDF0IBtbMTswOzEwMzsxMDh0sN+wG1sxOzA7MTQ1OzE0MXQgIBtbMTs0Nm0bWzA7MDsyMTc7MjM0dBtbMTsxMjg7MjU1OzI1NXSwG1s1OzQ1OzM1bRtbMDsyMzg7MTA0OzI1M3QbWzE7MjQ1OzE3NjsyNTV0sbAbWzA7MTszNW0bWzE7MjM4OzEwNDsyNTN02xtbNDY7MzZtG1swOzA7MjE3OzIzNHQbWzE7MTI4OzI1NTsyNTV0IBtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dLCw398bWzE7MDsxNDU7MTQxdCAgG1sxOzM1bRtbMTsyNTU7MTU1OzE1N3QgG1swOzM2bRtbMTswOzE0NTsxNDF0ICAgG1s0Nm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0sBtbNDBtG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLAbWzQwOzM1bRtbMTsyMzA7MDsyMjR0sRtbMzZtG1sxOzA7MTQ1OzE0MXQgIBtbNTs0NTszNW0bWzA7MjU1OzEwMjsxMDZ0G1sxOzE5MTswOzV0shtbMDsxOzM1bRtbMTsyNTU7MTAyOzEwNnQgG1swOzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSyG1s0MDszNW0bWzE7MjMwOzA7MjI0dLLfG1sxbRtbMTsyMzg7MTA0OzI1M3QgG1swOzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSyG1s0MG0bWzE7MDsxNDU7MTQxdCAgG1s0Nm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0sBtbNDA7MzVtG1sxOzIzMDswOzIyNHSxsrIbWzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dN8bWzE7MDsyMTc7MjM0dCAbWzE7MDsxMDM7MTA4dNwNChtbMTswOzIxNzsyMzR0ICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzA7MTAzOzEwOHSxG1s0MG3bG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0shtbNDU7MzVtG1swOzE2ODsyOzE1NXQbWzE7MjMwOzA7MjI0dLGwsBtbNDY7MzZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLEbWzQwbRtbMTswOzE0NTsxNDF0IBtbMTswOzEwMzsxMDh0IBtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLEbWzQ1OzM1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSwsLEbWzQ2OzM2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSxG1s0NTszNW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sLGxG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB0sRtbNDBtG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLEbWzQwOzM1bRtbMTsxNjg7MjsxNTV0sbGyshtbNDY7MzZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLIbWzQwbRtbMTswOzE0NTsxNDF0IBtbMTswOzEwMzsxMDh03RtbMTswOzE0NTsxNDF0ICAgIBtbMTs0Nm0bWzA7MDsyMTc7MjM0dBtbMTsxMjg7MjU1OzI1NXSxG1s1OzQ1OzM1bRtbMDsyMzg7MTA0OzI1M3QbWzE7MjQ1OzE3NjsyNTV03xtbMDsxOzM1bRtbMTsyMzg7MTA0OzI1M3Tb3xtbMDszNm0bWzE7MDsyMTc7MjM0dN8bWzE7MDsxNDU7MTQxdCAbWzE7MzVtG1sxOzI1NTsxMDI7MTA2dCAbWzA7MzVtG1sxOzE5MTswOzV03BtbMW0bWzE7MjU1OzEwMjsxMDZ03NwbWzQ1bRtbMDsxOTE7MDs1dBtbMTsyNTU7MTU1OzE1N3TcG1s0MG0bWzE7MjU1OzEwMjsxMDZ0398gG1swOzM2bRtbMTswOzE0NTsxNDF0ICAbWzE7MjsxOTA7MjAwdN8bWzE7MDsxNDU7MTQxdCAgG1sxOzI7MTkwOzIwMHTfG1sxOzA7MTQ1OzE0MXQgICAbWzE7MzVtG1sxOzI1NTsxNTU7MTU3dCAbWzA7NTs0NTszNW0bWzA7MjU1OzEwMjsxMDZ0G1sxOzE5MTswOzV0sRtbMDsxOzM1bRtbMTsyNTU7MTU1OzE1N3QgG1swOzM2bRtbMTswOzE0NTsxNDF0IBtbMTszNW0bWzE7MjU1OzEwMjsxMDZ03BtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzE7MzVtG1sxOzI1NTsxNTU7MTU3dCAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLEbWzQwOzM1bRtbMTsyMzA7MDsyMjR0srIbWzE7NDVtG1swOzIzMDswOzIyNHQbWzE7MjM4OzEwNDsyNTN0sBtbMDs0NjszNm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0shtbNDBtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHTc3xtbNDZtG1swOzA7MTQ1OzE0MXSyG1swbQ0KG1szNm0bWzE7MDsyMTc7MjM0dCAbWzE7MDsxMDM7MTA4dNvcG1sxOzA7MTQ1OzE0MXQgIBtbNDZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLAbWzQ1OzM1bRtbMDsxNjg7MjsxNTV0G1sxOzIzMDswOzIyNHSysbEbWzQ2OzM2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dLAbWzE7MDsxNDU7MTQxdCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0NTszNW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sbKyG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB02xtbNDBtG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR02xtbMTs0NTszNW0bWzA7MjMwOzA7MjI0dBtbMTsyMzg7MTA0OzI1M3Sw2xtbMDs0NTszNW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0shtbNDY7MzZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdLIbWzQwbRtbMTswOzE0NTsxNDF0ICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzI7MTkwOzIwMHSyG1s0MDszNW0bWzE7MTY4OzI7MTU1dLKyG1s0NW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sLAbWzE7NDY7MzZtG1swOzA7MjE3OzIzNHQbWzE7MTI4OzI1NTsyNTV0sBtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dN/csBtbMTswOzE0NTsxNDF0IBtbMTszNW0bWzE7MjU1OzE1NTsxNTd0ICAbWzA7MzVtG1sxOzE5MTswOzV03BtbMW0bWzE7MjU1OzEwMjsxMDZ03NwbWzA7NTs0NTszNW0bWzA7MjU1OzEwMjsxMDZ0G1sxOzE5MTswOzV03xtbMW0bWzE7MjU1OzE1NTsxNTd03NzfG1swOzE7MzVtG1sxOzI1NTsxNTU7MTU3dN8bWzE7MjU1OzEwMjsxMDZ0ICAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1s0Nm0bWzA7MjsxOTA7MjAwdBtbMTswOzIxNzsyMzR0IBtbMTs0MDszNW0bWzE7MjM4OzEwNDsyNTN0sBtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzE7MzVtG1sxOzI1NTsxMDI7MTA2dN/cIBtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzE7MzVtG1sxOzI1NTsxMDI7MTA2dCDbG1sxOzI1NTsxNTU7MTU3dCAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzM1bRtbMTsyNTU7MTU1OzE1N3QgG1swOzM1bRtbMTsxOTE7MDs1dNwbWzU7MTs0NW0bWzA7MjU1OzEwMjsxMDZ0G1sxOzI1NTsxNTU7MTU3dLIbWzA7MTszNW0bWzE7MjU1OzEwMjsxMDZ03yAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzM1bRtbMTsyNTU7MTAyOzEwNnTf3BtbMTsyNTU7MTU1OzE1N3QgG1swOzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSyG1s0MDszNW0bWzE7MjMwOzA7MjI0dNsbWzE7NDVtG1swOzIzMDswOzIyNHQbWzE7MjM4OzEwNDsyNTN03LEbWzQ2OzM2bRtbMDswOzIxNzsyMzR0G1sxOzEyODsyNTU7MjU1dLAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgIBtbMTswOzEwMzsxMDh03BtbNDZtG1swOzA7MTQ1OzE0MXSyG1swbQ0KG1szNm0bWzE7MDsyMTc7MjM0dCAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzA7MTAzOzEwOHSysrIbWzQwbRtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLEbWzE7NDU7MzVtG1swOzIzMDswOzIyNHQbWzE7MjM4OzEwNDsyNTN03BtbMDs0NTszNW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0srIbWzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dCAbWzE7MDsxNDU7MTQxdCAbWzE7NDZtG1swOzA7MjE3OzIzNHQbWzE7MTI4OzI1NTsyNTV0sBtbMDs0NTszNW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0shtbMW0bWzA7MjMwOzA7MjI0dBtbMTsyMzg7MTA0OzI1M3TcG1swOzM1bRtbMTsyMzA7MDsyMjR02xtbNDY7MzZtG1swOzA7MTQ1OzE0MXQbWzE7MjsxOTA7MjAwdNsbWzQwbRtbMTswOzE0NTsxNDF0IBtbMTs0Nm0bWzA7MDsyMTc7MjM0dBtbMTsxMjg7MjU1OzI1NXSwG1s1OzQ1OzM1bRtbMDsyMzg7MTA0OzI1M3QbWzE7MjQ1OzE3NjsyNTV03N8bWzA7MzVtG1sxOzIzMDswOzIyNHTfG1s0NjszNm0bWzA7MDsxNDU7MTQxdBtbMTsyOzE5MDsyMDB02xtbNDBtG1sxOzA7MjE3OzIzNHTc3xtbMzVtG1sxOzIzMDswOzIyNHTcG1s0NW0bWzA7MTY4OzI7MTU1dLIbWzFtG1swOzIzMDswOzIyNHQbWzE7MjM4OzEwNDsyNTN03BtbMDs0NTszNW0bWzA7MTY4OzI7MTU1dBtbMTsyMzA7MDsyMjR0sbEbWzE7NDY7MzZtG1swOzA7MjE3OzIzNHQbWzE7MTI4OzI1NTsyNTV0shtbMG0bWzVDG1szNW0bWzE7MTkxOzA7NXTfG1sxbRtbMTsyNTU7MTAyOzEwNnTf3xtbMTsyNTU7MTU1OzE1N3Tf3yAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgIBtbMzVtG1sxOzE5MTswOzV02xtbNTsxOzQ1bRtbMDsyNTU7MTAyOzEwNnQbWzE7MjU1OzE1NTsxNTd0IBtbMDszNm0bWzE7MDsxNDU7MTQxdCAgIBtbMTs0Nm0bWzA7MDsyMTc7MjM0dBtbMTsxMjg7MjU1OzI1NXSwG1s0MDszNW0bWzE7MjM4OzEwNDsyNTN0srAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1szNW0bWzE7MTkxOzA7NXTfG1s1OzQ1bRtbMDsyNTU7MTAyOzEwNnTcG1swOzE7MzVtG1sxOzI1NTsxMDI7MTA2dNwbWzU7NDVtG1swOzI1NTsxMDI7MTA2dBtbMTsyNTU7MTU1OzE1N3TcG1swOzE7MzVtG1sxOzI1NTsxMDI7MTA2dN0gG1swOzM2bRtbMTswOzE0NTsxNDF0IBtbMzVtG1sxOzE5MTswOzV03xtbMW0bWzE7MjU1OzEwMjsxMDZ0IBtbMDszNW0bWzE7MTkxOzA7NXTbG1s1OzE7NDVtG1swOzI1NTsxMDI7MTA2dBtbMTsyNTU7MTU1OzE1N3SxG1swOzE7MzVtG1sxOzI1NTsxNTU7MTU3dCAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzM1bRtbMTsyNTU7MTAyOzEwNnQgG1s1OzQ1bRtbMDsyNTU7MTAyOzEwNnQbWzE7MjU1OzE1NTsxNTd02xtbMDsxOzM1bRtbMTsyNTU7MTAyOzEwNnTdG1swOzM2bRtbMTswOzE0NTsxNDF0IBtbMTs0NTszNW0bWzA7MjMwOzA7MjI0dBtbMTsyMzg7MTA0OzI1M3TbG1s1bRtbMDsyMzg7MTA0OzI1M3QbWzE7MjQ1OzE3NjsyNTV03BtbMDsxOzQ1OzM1bRtbMDsyMzA7MDsyMjR0G1sxOzIzODsxMDQ7MjUzdLIbWzQ2OzM2bRtbMDswOzIxNzsyMzR0G1sxOzEyODsyNTU7MjU1dLEbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHTfG1s0Nm0bWzA7MDsxNDU7MTQxdLGwG1swbQ0KG1szNm0bWzE7MDsyMTc7MjM0dCAbWzE7MDsxMDM7MTA4dNzfG1s0Nm0bWzA7MDsxNDU7MTQxdNwbWzQwbRtbMTswOzE0NTsxNDF0IBtbNDZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLIbWzE7NDU7MzVtG1swOzIzMDswOzIyNHQbWzE7MjQ1OzE3NjsyNTV03BtbMTsyMzg7MTA0OzI1M3TfG1sxOzI0NTsxNzY7MjU1dN8bWzA7NDY7MzZtG1swOzI7MTkwOzIwMHQbWzE7MDsyMTc7MjM0dLEbWzQwbRtbMTswOzE0NTsxNDF0IBtbMTswOzEwMzsxMDh0sBtbMTswOzE0NTsxNDF0IBtbMTs0Nm0bWzA7MDsyMTc7MjM0dBtbMTsxMjg7MjU1OzI1NXSyG1s1OzQ1OzM1bRtbMDsyMzg7MTA0OzI1M3QbWzE7MjQ1OzE3NjsyNTV039vcG1swOzQ2OzM2bRtbMDsyOzE5MDsyMDB0G1sxOzA7MjE3OzIzNHSwG1s0MG0bWzE7MDsxNDU7MTQxdCAbWzE7NDZtG1swOzA7MjE3OzIzNHQbWzE7MTI4OzI1NTsyNTV0sRtbNTs0NTszNW0bWzA7MjM4OzEwNDsyNTN0G1sxOzI0NTsxNzY7MjU1dNwbWzA7MzVtG1sxOzIzMDswOzIyNHTcG1szNm0bWzE7MDsyMTc7MjM0dN8bWzM1bRtbMTsyMzA7MDsyMjR03NsbWzE7NDVtG1swOzIzMDswOzIyNHQbWzE7MjM4OzEwNDsyNTN03BtbNW0bWzA7MjM4OzEwNDsyNTN0G1sxOzI0NTsxNzY7MjU1dNzc3xtbMDszNW0bWzE7MjMwOzA7MjI0dN8bWzE7MzZtG1sxOzEyODsyNTU7MjU1dNsbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgIBtbMTswOzEwMzsxMDh0sLCw3BtbMTswOzE0NTsxNDF0IBtbMTs0Nm0bWzA7MDsyMTc7MjM0dBtbMTsxMjg7MjU1OzI1NXSyG1s0MDszNW0bWzE7MjM4OzEwNDsyNTN02xtbNTs0NW0bWzA7MjM4OzEwNDsyNTN0G1sxOzI0NTsxNzY7MjU1dNzfG1swOzE7NDY7MzZtG1swOzA7MjE3OzIzNHQbWzE7MTI4OzI1NTsyNTV0sBtbMDszNm0bWzE7MDsxNDU7MTQxdCAgIBtbMzVtG1sxOzE5MTswOzV02xtbNTsxOzQ1bRtbMDsyNTU7MTAyOzEwNnQbWzE7MjU1OzE1NTsxNTd0sBtbMDszNm0bWzE7MDsxNDU7MTQxdCAgG1sxOzQ2bRtbMDswOzIxNzsyMzR0G1sxOzEyODsyNTU7MjU1dLEbWzQwOzM1bRtbMTsyMzg7MTA0OzI1M3Tb2xtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzE7MzVtG1sxOzI1NTsxNTU7MTU3dCAbWzE7MjU1OzEwMjsxMDZ0IBtbMDs1OzQ1OzM1bRtbMDsyNTU7MTAyOzEwNnQbWzE7MTkxOzA7NXSxG1swOzE7NDU7MzVtG1swOzE5MTswOzV0G1sxOzI1NTsxNTU7MTU3dLIbWzQwbRtbMTsyNTU7MTAyOzEwNnTcICAbWzE7MjU1OzE1NTsxNTd0IBtbMDszNm0bWzE7MDsxNDU7MTQxdCAbWzM1bRtbMTsxOTE7MDs1dCCyG1s1OzE7NDVtG1swOzI1NTsxMDI7MTA2dBtbMTsyNTU7MTU1OzE1N3TfG1swOzE7MzVtG1sxOzI1NTsxNTU7MTU3dNwbWzU7NDVtG1swOzI1NTsxMDI7MTA2dN8bWzA7MTszNW0bWzE7MjU1OzEwMjsxMDZ03xtbMTsyNTU7MTU1OzE1N3QgG1swOzM2bRtbMTswOzE0NTsxNDF0IBtbMzVtG1sxOzIzMDswOzIyNHTfG1sxOzQ1bRtbMDsyMzA7MDsyMjR0G1sxOzIzODsxMDQ7MjUzdN8bWzA7MzVtG1sxOzIzMDswOzIyNHTbG1sxOzQ2OzM2bRtbMDswOzIxNzsyMzR0G1sxOzEyODsyNTU7MjU1dLIbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHTbG1s0Nm0bWzA7MDsxNDU7MTQxdNwbWzQwbd8NChtbMTswOzIxNzsyMzR0IBtbNDZtG1swOzA7MTQ1OzE0MXQbWzE7MDsxMDM7MTA4dNwbWzQwbd/cG1sxOzA7MTQ1OzE0MXQgG1sxOzQ2bRtbMDswOzIxNzsyMzR0G1sxOzEyODsyNTU7MjU1dLIbWzQwbdzc3BtbNDZtG1swOzA7MjE3OzIzNHSxG1swOzM2bRtbMTswOzE0NTsxNDF0IBtbMTswOzEwMzsxMDh0IBtbMTswOzE0NTsxNDF0IBtbMW0bWzE7MTI4OzI1NTsyNTV029zc3NsbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzQ2bRtbMDswOzIxNzsyMzR0G1sxOzEyODsyNTU7MjU1dLIbWzQwbdzc3Nzc3Nzc398bWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHTcG1szN20bWzZDG1sxOzM2bRtbMTsxMjg7MjU1OzI1NXTb3NzcG1s0Nm0bWzA7MDsyMTc7MjM0dLIbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1szNW0bWzE7MTkxOzA7NXTcG1szNm0bWzE7MDsxNDU7MTQxdCAgG1szNW0bWzE7MTkxOzA7NXTbG1s1OzE7NDVtG1swOzI1NTsxMDI7MTA2dBtbMTsyNTU7MTU1OzE1N3SxG1swOzM2bRtbMTswOzE0NTsxNDF0IBtbMTswOzIxNzsyMzR03xtbMW0bWzE7MTI4OzI1NTsyNTV03NwbWzM1bRtbMTsyNTU7MTU1OzE1N3QgG1sxOzI1NTsxMDI7MTA2dCAbWzA7MzVtG1sxOzE5MTswOzV03N8bWzFtG1sxOzI1NTsxNTU7MTU3dCAgG1sxOzI1NTsxMDI7MTA2dN/cIBtbMDszNm0bWzE7MDsxNDU7MTQxdCAgIBtbMzVtG1sxOzE5MTswOzV03xtbMW0bWzE7MjU1OzEwMjsxMDZ03yAbWzA7MzZtG1sxOzA7MTQ1OzE0MXQgIBtbMW0bWzE7MTI4OzI1NTsyNTV039/c3N8bWzA7MzZtG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHTc3xtbNDZtG1swOzA7MTQ1OzE0MXSyG1swbQ0KG1szNm0bWzE7MDsyMTc7MjM0dCAbWzE7MDsxMDM7MTA4dNwbWzQ2bRtbMDswOzE0NTsxNDF037AbWzQwbbCwsBtbMTswOzIxNzsyMzR0ICAgIBtbMTswOzEwMzsxMDh0IBtbMzdtG1sxNkMbWzM2bRtbMTswOzEwMzsxMDh03NzfG1sxOzA7MjE3OzIzNHQgIBtbMTswOzEwMzsxMDh0sLCwG1szN20bWzhDG1szNW0bWzE7MTkxOzA7NXTeG1s1OzQ1bRtbMDsyNTU7MTAyOzEwNnSyG1swOzM2bRtbMTswOzIxNzsyMzR0ICAbWzM1bRtbMTsxOTE7MDs1dLIbWzFtG1sxOzI1NTsxMDI7MTA2dNsbWzA7MzZtG1sxOzA7MjE3OzIzNHQgIBtbMzVtG1sxOzE5MTswOzV03BtbMW0bWzE7MjU1OzE1NTsxNTd0IBtbMTsyNTU7MTAyOzEwNnQgG1swOzM1bRtbMTsxOTE7MDs1dN8bWzFtG1sxOzI1NTsxMDI7MTA2dCAbWzA7MzZtG1sxOzA7MjE3OzIzNHQgG1szNW0bWzE7MTkxOzA7NXQgIBtbMzZtG1sxOzA7MjE3OzIzNHQgG1szN20bWzEzQxtbMzZtG1sxOzA7MTAzOzEwOHSwsBtbNDZtG1swOzA7MTQ1OzE0MXSxsLEbWzBtDQobWzM2bRtbMTswOzIxNzsyMzR0ICAbWzQ2bRtbMDswOzE0NTsxNDF0G1sxOzA7MTAzOzEwOHSy3BtbNDBt39/f398bWzE7MDsxNDU7MTQxdCAgG1sxOzA7MTAzOzEwOHTf398bWzE7MDsxNDU7MTQxdCAgG1sxOzA7MTAzOzEwOHTfG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHTfG1sxOzA7MTQ1OzE0MXQgICAbWzE7MDsxMDM7MTA4dN8bWzE7MDsxNDU7MTQxdCAbWzE7MDsxMDM7MTA4dN/f39/fG1szN20bWzEwQxtbMzVtG1sxOzE5MTswOzV03NwbWzM2bRtbMTswOzE0NTsxNDF0IBtbMzVtG1sxOzE5MTswOzV03NzfG1s1OzQ1bRtbMDsyNTU7MTAyOzEwNnTcG1swOzM1bRtbMTsxOTE7MDs1dLIgILIbWzFtG1sxOzI1NTsxMDI7MTA2dN0bWzA7MzVtG1sxOzE5MTswOzV0IBtbNTs0NW0bWzA7MjU1OzEwMjsxMDZ03xtbMDszNW0bWzE7MTkxOzA7NXSy3Nwg3NwbWzM3bRtbOEMbWzM2bRtbMTswOzEwMzsxMDh0WklJ39/fG1sxOzA7MTQ1OzE0MXQgG1sxOzA7MTAzOzEwOHTf398bWzQ2bRtbMDswOzE0NTsxNDF03BtbNDBt2w0KDQobWzM3bSAgG1sxbVJldHJvVHh0IC0bWzBtIBtbMTszNm1UaGUgYXBwIHRoYXQgbGV0cyB5b3UgdmlldyB3b3JrcyBvZiBBTlNJIGFydCwgQVNDSUkgYW5kIE5GTyB0ZXh0DQobWzBtICAbWzE7MzZtaW4tYnJvd3NlciBhcyBIVE1MLiBMR1BMIGFuZCBhdmFpbGFibGUgZm9yIENocm9tZSwgRmlyZWZveC4NChtbMG0gIBtbMTszNW1odHRwOi8vcmV0cm90eHQuY29tDQobWzBtICAbWzM2bUxvZ28gYnkgWmV1cyBJSSBbQmxvY2t0cm9uaWNzXSBbRlVFTF0NChpDT01OVFRoZSBhcHAgdGhhdCBsZXRzIHlvdSB2aWV3IHdvcmtzIG9mIEFOU0kgYXJ0LCBBU0NJSSBhbmQgTkZPIHRleHQgaW4tYnJvd3NlciBhcyBIVE1MLiBMR1BMIGFuZCBhdmFpbGFibGUgZm9yIENocm9tZSwgRmlyZWZveC4gICAgU0FVQ0UwMFJldHJvVHh0IGxvZ28gICAgICAgICAgICAgICAgICAgICAgWmV1cyBJSSAgICAgICAgICAgICBCbG9ja3Ryb25pY3MsIEZVRUwgIDIwMTcwNzMwGmYAAAEBUAAfAAAAAAACE0lCTSBWR0EAAAAAAAAAAAAAAAAAAAANCg=="

// EncodeASCII encodes the content of ?.asc to base64 for use as LogoASCII.
func EncodeASCII() (result string) {
	d, err := ReadLine(asciiFile(), "dos")
	if err != nil {
		log.Fatal(err)
	}
	return Base64Encode(d)
}

// EncodeANSI encodes the content of ZII-RTXT.ans to base64 for use as LogoANSI.
func EncodeANSI() (result string) {
	d, err := ReadLine(ansiFile(), "dos")
	if err != nil {
		log.Fatal(err)
	}
	return Base64Encode(d)
}
