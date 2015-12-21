/*
Go Language Raspberry Pi Interface
(c) Copyright David Thorpe 2016
All Rights Reserved

For Licensing and Usage information, please see LICENSE.md
*/
package dispmanx

const (
	/* Bottom 2 bits sets the alpha mode */
	DISPMANX_FLAGS_ALPHA_FROM_SOURCE       = 0
	DISPMANX_FLAGS_ALPHA_FIXED_ALL_PIXELS  = 1
	DISPMANX_FLAGS_ALPHA_FIXED_NON_ZERO    = 2
	DISPMANX_FLAGS_ALPHA_FIXED_EXCEED_0X07 = 3
	DISPMANX_FLAGS_ALPHA_PREMULT           = 1 << 16
	DISPMANX_FLAGS_ALPHA_MIX               = 1 << 17

	/* Bottom 2 bits sets the orientation */
	DISPMANX_NO_ROTATE  = 0
	DISPMANX_ROTATE_90  = 1
	DISPMANX_ROTATE_180 = 2
	DISPMANX_ROTATE_270 = 3
	DISPMANX_FLIP_HRIZ  = 1 << 16
	DISPMANX_FLIP_VERT  = 1 << 17

	DISPMANX_PROTECTION_MAX  = 0x0f
	DISPMANX_PROTECTION_NONE = 0
	DISPMANX_PROTECTION_HDCP = 11 // Derived from the WM DRM levels, 101-300

)
