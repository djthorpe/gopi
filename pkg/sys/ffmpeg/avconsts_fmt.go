package ffmpeg

const (
	AV_PIX_FMT_NONE           AVPixelFormat = iota
	AV_PIX_FMT_YUV420P                      // planar YUV 4:2:0, 12bpp, (1 Cr & Cb sample per 2x2 Y samples)
	AV_PIX_FMT_YUYV422                      // packed YUV 4:2:2, 16bpp, Y0 Cb Y1 Cr
	AV_PIX_FMT_RGB24                        // packed RGB 8:8:8, 24bpp, RGBRGB...
	AV_PIX_FMT_BGR24                        // packed RGB 8:8:8, 24bpp, BGRBGR...
	AV_PIX_FMT_YUV422P                      // planar YUV 4:2:2, 16bpp, (1 Cr & Cb sample per 2x1 Y samples)
	AV_PIX_FMT_YUV444P                      // planar YUV 4:4:4, 24bpp, (1 Cr & Cb sample per 1x1 Y samples)
	AV_PIX_FMT_YUV410P                      // planar YUV 4:1:0, 9bpp, (1 Cr & Cb sample per 4x4 Y samples)
	AV_PIX_FMT_YUV411P                      // planar YUV 4:1:1, 12bpp, (1 Cr & Cb sample per 4x1 Y samples)
	AV_PIX_FMT_GRAY8                        // 8bpp.
	AV_PIX_FMT_MONOWHITE                    // 1bpp, 0 is white, 1 is black, in each byte pixels are ordered from the msb to the lsb.
	AV_PIX_FMT_MONOBLACK                    // 1bpp, 0 is black, 1 is white, in each byte pixels are ordered from the msb to the lsb.
	AV_PIX_FMT_PAL8                         // 8 bits with AV_PIX_FMT_RGB32 palette
	AV_PIX_FMT_YUVJ420P                     // planar YUV 4:2:0, 12bpp, full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV420P and setting color_range
	AV_PIX_FMT_YUVJ422P                     // planar YUV 4:2:2, 16bpp, full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV422P and setting color_range
	AV_PIX_FMT_YUVJ444P                     // planar YUV 4:4:4, 24bpp, full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV444P and setting color_range
	AV_PIX_FMT_UYVY422                      // packed YUV 4:2:2, 16bpp, Cb Y0 Cr Y1
	AV_PIX_FMT_UYYVYY411                    // packed YUV 4:1:1, 12bpp, Cb Y0 Y1 Cr Y2 Y3
	AV_PIX_FMT_BGR8                         // packed RGB 3:3:2, 8bpp, (msb)2B 3G 3R(lsb)
	AV_PIX_FMT_BGR4                         // packed RGB 1:2:1 bitstream, 4bpp, (msb)1B 2G 1R(lsb), a byte contains two pixels, the first pixel in the byte is the one composed by the 4 msb bits
	AV_PIX_FMT_BGR4_BYTE                    // packed RGB 1:2:1, 8bpp, (msb)1B 2G 1R(lsb)
	AV_PIX_FMT_RGB8                         // packed RGB 3:3:2, 8bpp, (msb)2R 3G 3B(lsb)
	AV_PIX_FMT_RGB4                         // packed RGB 1:2:1 bitstream, 4bpp, (msb)1R 2G 1B(lsb), a byte contains two pixels, the first pixel in the byte is the one composed by the 4 msb bits
	AV_PIX_FMT_RGB4_BYTE                    // packed RGB 1:2:1, 8bpp, (msb)1R 2G 1B(lsb)
	AV_PIX_FMT_NV12                         // planar YUV 4:2:0, 12bpp, 1 plane for Y and 1 plane for the UV components, which are interleaved (first byte U and the following byte V)
	AV_PIX_FMT_NV21                         // as above, but U and V bytes are swapped
	AV_PIX_FMT_ARGB                         // packed ARGB 8:8:8:8, 32bpp, ARGBARGB...
	AV_PIX_FMT_RGBA                         // packed RGBA 8:8:8:8, 32bpp, RGBARGBA...
	AV_PIX_FMT_ABGR                         // packed ABGR 8:8:8:8, 32bpp, ABGRABGR...
	AV_PIX_FMT_BGRA                         // packed BGRA 8:8:8:8, 32bpp, BGRABGRA...
	AV_PIX_FMT_GRAY16BE                     // 16bpp, big-endian.
	AV_PIX_FMT_GRAY16LE                     // 16bpp, little-endian.
	AV_PIX_FMT_YUV440P                      // planar YUV 4:4:0 (1 Cr & Cb sample per 1x2 Y samples)
	AV_PIX_FMT_YUVJ440P                     // planar YUV 4:4:0 full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV440P and setting color_range
	AV_PIX_FMT_YUVA420P                     // planar YUV 4:2:0, 20bpp, (1 Cr & Cb sample per 2x2 Y & A samples)
	AV_PIX_FMT_RGB48BE                      // packed RGB 16:16:16, 48bpp, 16R, 16G, 16B, the 2-byte value for each R/G/B component is stored as big-endian
	AV_PIX_FMT_RGB48LE                      // packed RGB 16:16:16, 48bpp, 16R, 16G, 16B, the 2-byte value for each R/G/B component is stored as little-endian
	AV_PIX_FMT_RGB565BE                     // packed RGB 5:6:5, 16bpp, (msb) 5R 6G 5B(lsb), big-endian
	AV_PIX_FMT_RGB565LE                     // packed RGB 5:6:5, 16bpp, (msb) 5R 6G 5B(lsb), little-endian
	AV_PIX_FMT_RGB555BE                     // packed RGB 5:5:5, 16bpp, (msb)1X 5R 5G 5B(lsb), big-endian , X=unused/undefined
	AV_PIX_FMT_RGB555LE                     // packed RGB 5:5:5, 16bpp, (msb)1X 5R 5G 5B(lsb), little-endian, X=unused/undefined
	AV_PIX_FMT_BGR565BE                     // packed BGR 5:6:5, 16bpp, (msb) 5B 6G 5R(lsb), big-endian
	AV_PIX_FMT_BGR565LE                     // packed BGR 5:6:5, 16bpp, (msb) 5B 6G 5R(lsb), little-endian
	AV_PIX_FMT_BGR555BE                     // packed BGR 5:5:5, 16bpp, (msb)1X 5B 5G 5R(lsb), big-endian , X=unused/undefined
	AV_PIX_FMT_BGR555LE                     // packed BGR 5:5:5, 16bpp, (msb)1X 5B 5G 5R(lsb), little-endian, X=unused/undefined
	AV_PIX_FMT_VAAPI_MOCO                   // HW acceleration through VA API at motion compensation entry-point, Picture.data[3] contains a vaapi_render_state struct which contains macroblocks as well as various fields extracted from headers.
	AV_PIX_FMT_VAAPI_IDCT                   // HW acceleration through VA API at IDCT entry-point, Picture.data[3] contains a vaapi_render_state struct which contains fields extracted from headers.
	AV_PIX_FMT_VAAPI_VLD                    // HW decoding through VA API, Picture.data[3] contains a VASurfaceID.
	AV_PIX_FMT_VAAPI                        //
	AV_PIX_FMT_YUV420P16LE                  // planar YUV 4:2:0, 24bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	AV_PIX_FMT_YUV420P16BE                  // planar YUV 4:2:0, 24bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	AV_PIX_FMT_YUV422P16LE                  // planar YUV 4:2:2, 32bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	AV_PIX_FMT_YUV422P16BE                  // planar YUV 4:2:2, 32bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	AV_PIX_FMT_YUV444P16LE                  // planar YUV 4:4:4, 48bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	AV_PIX_FMT_YUV444P16BE                  // planar YUV 4:4:4, 48bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	AV_PIX_FMT_DXVA2_VLD                    // HW decoding through DXVA2, Picture.data[3] contains a LPDIRECT3DSURFACE9 pointer.
	AV_PIX_FMT_RGB444LE                     // packed RGB 4:4:4, 16bpp, (msb)4X 4R 4G 4B(lsb), little-endian, X=unused/undefined
	AV_PIX_FMT_RGB444BE                     // packed RGB 4:4:4, 16bpp, (msb)4X 4R 4G 4B(lsb), big-endian, X=unused/undefined
	AV_PIX_FMT_BGR444LE                     // packed BGR 4:4:4, 16bpp, (msb)4X 4B 4G 4R(lsb), little-endian, X=unused/undefined
	AV_PIX_FMT_BGR444BE                     // packed BGR 4:4:4, 16bpp, (msb)4X 4B 4G 4R(lsb), big-endian, X=unused/undefined
	AV_PIX_FMT_YA8                          // 8 bits gray, 8 bits alpha
	AV_PIX_FMT_Y400A                        // alias for AV_PIX_FMT_YA8
	AV_PIX_FMT_GRAY8A                       // alias for AV_PIX_FMT_YA8
	AV_PIX_FMT_BGR48BE                      // packed RGB 16:16:16, 48bpp, 16B, 16G, 16R, the 2-byte value for each R/G/B component is stored as big-endian
	AV_PIX_FMT_BGR48LE                      // packed RGB 16:16:16, 48bpp, 16B, 16G, 16R, the 2-byte value for each R/G/B component is stored as little-endian
	AV_PIX_FMT_YUV420P9BE                   // The following 12 formats have the disadvantage of needing 1 format for each bit depth.
	AV_PIX_FMT_YUV420P9LE                   // planar YUV 4:2:0, 13.5bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	AV_PIX_FMT_YUV420P10BE                  // planar YUV 4:2:0, 15bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	AV_PIX_FMT_YUV420P10LE                  // planar YUV 4:2:0, 15bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	AV_PIX_FMT_YUV422P10BE                  // planar YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	AV_PIX_FMT_YUV422P10LE                  // planar YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	AV_PIX_FMT_YUV444P9BE                   // planar YUV 4:4:4, 27bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	AV_PIX_FMT_YUV444P9LE                   // planar YUV 4:4:4, 27bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	AV_PIX_FMT_YUV444P10BE                  // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	AV_PIX_FMT_YUV444P10LE                  // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	AV_PIX_FMT_YUV422P9BE                   // planar YUV 4:2:2, 18bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	AV_PIX_FMT_YUV422P9LE                   // planar YUV 4:2:2, 18bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	AV_PIX_FMT_GBRP                         //
	AV_PIX_FMT_GBR24P                       // planar GBR 4:4:4 24bpp
	AV_PIX_FMT_GBRP9BE                      // planar GBR 4:4:4 27bpp, big-endian
	AV_PIX_FMT_GBRP9LE                      // planar GBR 4:4:4 27bpp, little-endian
	AV_PIX_FMT_GBRP10BE                     // planar GBR 4:4:4 30bpp, big-endian
	AV_PIX_FMT_GBRP10LE                     // planar GBR 4:4:4 30bpp, little-endian
	AV_PIX_FMT_GBRP16BE                     // planar GBR 4:4:4 48bpp, big-endian
	AV_PIX_FMT_GBRP16LE                     // planar GBR 4:4:4 48bpp, little-endian
	AV_PIX_FMT_YUVA422P                     // planar YUV 4:2:2 24bpp, (1 Cr & Cb sample per 2x1 Y & A samples)
	AV_PIX_FMT_YUVA444P                     // planar YUV 4:4:4 32bpp, (1 Cr & Cb sample per 1x1 Y & A samples)
	AV_PIX_FMT_YUVA420P9BE                  // planar YUV 4:2:0 22.5bpp, (1 Cr & Cb sample per 2x2 Y & A samples), big-endian
	AV_PIX_FMT_YUVA420P9LE                  // planar YUV 4:2:0 22.5bpp, (1 Cr & Cb sample per 2x2 Y & A samples), little-endian
	AV_PIX_FMT_YUVA422P9BE                  // planar YUV 4:2:2 27bpp, (1 Cr & Cb sample per 2x1 Y & A samples), big-endian
	AV_PIX_FMT_YUVA422P9LE                  // planar YUV 4:2:2 27bpp, (1 Cr & Cb sample per 2x1 Y & A samples), little-endian
	AV_PIX_FMT_YUVA444P9BE                  // planar YUV 4:4:4 36bpp, (1 Cr & Cb sample per 1x1 Y & A samples), big-endian
	AV_PIX_FMT_YUVA444P9LE                  // planar YUV 4:4:4 36bpp, (1 Cr & Cb sample per 1x1 Y & A samples), little-endian
	AV_PIX_FMT_YUVA420P10BE                 // planar YUV 4:2:0 25bpp, (1 Cr & Cb sample per 2x2 Y & A samples, big-endian)
	AV_PIX_FMT_YUVA420P10LE                 // planar YUV 4:2:0 25bpp, (1 Cr & Cb sample per 2x2 Y & A samples, little-endian)
	AV_PIX_FMT_YUVA422P10BE                 // planar YUV 4:2:2 30bpp, (1 Cr & Cb sample per 2x1 Y & A samples, big-endian)
	AV_PIX_FMT_YUVA422P10LE                 // planar YUV 4:2:2 30bpp, (1 Cr & Cb sample per 2x1 Y & A samples, little-endian)
	AV_PIX_FMT_YUVA444P10BE                 // planar YUV 4:4:4 40bpp, (1 Cr & Cb sample per 1x1 Y & A samples, big-endian)
	AV_PIX_FMT_YUVA444P10LE                 // planar YUV 4:4:4 40bpp, (1 Cr & Cb sample per 1x1 Y & A samples, little-endian)
	AV_PIX_FMT_YUVA420P16BE                 // planar YUV 4:2:0 40bpp, (1 Cr & Cb sample per 2x2 Y & A samples, big-endian)
	AV_PIX_FMT_YUVA420P16LE                 // planar YUV 4:2:0 40bpp, (1 Cr & Cb sample per 2x2 Y & A samples, little-endian)
	AV_PIX_FMT_YUVA422P16BE                 // planar YUV 4:2:2 48bpp, (1 Cr & Cb sample per 2x1 Y & A samples, big-endian)
	AV_PIX_FMT_YUVA422P16LE                 // planar YUV 4:2:2 48bpp, (1 Cr & Cb sample per 2x1 Y & A samples, little-endian)
	AV_PIX_FMT_YUVA444P16BE                 // planar YUV 4:4:4 64bpp, (1 Cr & Cb sample per 1x1 Y & A samples, big-endian)
	AV_PIX_FMT_YUVA444P16LE                 // planar YUV 4:4:4 64bpp, (1 Cr & Cb sample per 1x1 Y & A samples, little-endian)
	AV_PIX_FMT_VDPAU                        // HW acceleration through VDPAU, Picture.data[3] contains a VdpVideoSurface.
	AV_PIX_FMT_XYZ12LE                      // packed XYZ 4:4:4, 36 bpp, (msb) 12X, 12Y, 12Z (lsb), the 2-byte value for each X/Y/Z is stored as little-endian, the 4 lower bits are set to 0
	AV_PIX_FMT_XYZ12BE                      // packed XYZ 4:4:4, 36 bpp, (msb) 12X, 12Y, 12Z (lsb), the 2-byte value for each X/Y/Z is stored as big-endian, the 4 lower bits are set to 0
	AV_PIX_FMT_NV16                         // interleaved chroma YUV 4:2:2, 16bpp, (1 Cr & Cb sample per 2x1 Y samples)
	AV_PIX_FMT_NV20LE                       // interleaved chroma YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	AV_PIX_FMT_NV20BE                       // interleaved chroma YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	AV_PIX_FMT_RGBA64BE                     // packed RGBA 16:16:16:16, 64bpp, 16R, 16G, 16B, 16A, the 2-byte value for each R/G/B/A component is stored as big-endian
	AV_PIX_FMT_RGBA64LE                     // packed RGBA 16:16:16:16, 64bpp, 16R, 16G, 16B, 16A, the 2-byte value for each R/G/B/A component is stored as little-endian
	AV_PIX_FMT_BGRA64BE                     // packed RGBA 16:16:16:16, 64bpp, 16B, 16G, 16R, 16A, the 2-byte value for each R/G/B/A component is stored as big-endian
	AV_PIX_FMT_BGRA64LE                     // packed RGBA 16:16:16:16, 64bpp, 16B, 16G, 16R, 16A, the 2-byte value for each R/G/B/A component is stored as little-endian
	AV_PIX_FMT_YVYU422                      // packed YUV 4:2:2, 16bpp, Y0 Cr Y1 Cb
	AV_PIX_FMT_YA16BE                       // 16 bits gray, 16 bits alpha (big-endian)
	AV_PIX_FMT_YA16LE                       // 16 bits gray, 16 bits alpha (little-endian)
	AV_PIX_FMT_GBRAP                        // planar GBRA 4:4:4:4 32bpp
	AV_PIX_FMT_GBRAP16BE                    // planar GBRA 4:4:4:4 64bpp, big-endian
	AV_PIX_FMT_GBRAP16LE                    // planar GBRA 4:4:4:4 64bpp, little-endian
	AV_PIX_FMT_QSV                          // HW acceleration through QSV, data[3] contains a pointer to the mfxFrameSurface1 structure.
	AV_PIX_FMT_MMAL                         // HW acceleration though MMAL, data[3] contains a pointer to the MMAL_BUFFER_HEADER_T structure.
	AV_PIX_FMT_D3D11VA_VLD                  // HW decoding through Direct3D11 via old API, Picture.data[3] contains a ID3D11VideoDecoderOutputView pointer.
	AV_PIX_FMT_CUDA                         // HW acceleration through CUDA.
	AV_PIX_FMT_0RGB                         // packed RGB 8:8:8, 32bpp, XRGBXRGB... X=unused/undefined
	AV_PIX_FMT_RGB0                         // packed RGB 8:8:8, 32bpp, RGBXRGBX... X=unused/undefined
	AV_PIX_FMT_0BGR                         // packed BGR 8:8:8, 32bpp, XBGRXBGR... X=unused/undefined
	AV_PIX_FMT_BGR0                         // packed BGR 8:8:8, 32bpp, BGRXBGRX... X=unused/undefined
	AV_PIX_FMT_YUV420P12BE                  // planar YUV 4:2:0,18bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	AV_PIX_FMT_YUV420P12LE                  // planar YUV 4:2:0,18bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	AV_PIX_FMT_YUV420P14BE                  // planar YUV 4:2:0,21bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	AV_PIX_FMT_YUV420P14LE                  // planar YUV 4:2:0,21bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	AV_PIX_FMT_YUV422P12BE                  // planar YUV 4:2:2,24bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	AV_PIX_FMT_YUV422P12LE                  // planar YUV 4:2:2,24bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	AV_PIX_FMT_YUV422P14BE                  // planar YUV 4:2:2,28bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	AV_PIX_FMT_YUV422P14LE                  // planar YUV 4:2:2,28bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	AV_PIX_FMT_YUV444P12BE                  // planar YUV 4:4:4,36bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	AV_PIX_FMT_YUV444P12LE                  // planar YUV 4:4:4,36bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	AV_PIX_FMT_YUV444P14BE                  // planar YUV 4:4:4,42bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	AV_PIX_FMT_YUV444P14LE                  // planar YUV 4:4:4,42bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	AV_PIX_FMT_GBRP12BE                     // planar GBR 4:4:4 36bpp, big-endian
	AV_PIX_FMT_GBRP12LE                     // planar GBR 4:4:4 36bpp, little-endian
	AV_PIX_FMT_GBRP14BE                     // planar GBR 4:4:4 42bpp, big-endian
	AV_PIX_FMT_GBRP14LE                     // planar GBR 4:4:4 42bpp, little-endian
	AV_PIX_FMT_YUVJ411P                     // planar YUV 4:1:1, 12bpp, (1 Cr & Cb sample per 4x1 Y samples) full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV411P and setting color_range
	AV_PIX_FMT_BAYER_BGGR8                  // bayer, BGBG..(odd line), GRGR..(even line), 8-bit samples
	AV_PIX_FMT_BAYER_RGGB8                  // bayer, RGRG..(odd line), GBGB..(even line), 8-bit samples
	AV_PIX_FMT_BAYER_GBRG8                  // bayer, GBGB..(odd line), RGRG..(even line), 8-bit samples
	AV_PIX_FMT_BAYER_GRBG8                  // bayer, GRGR..(odd line), BGBG..(even line), 8-bit samples
	AV_PIX_FMT_BAYER_BGGR16LE               // bayer, BGBG..(odd line), GRGR..(even line), 16-bit samples, little-endian
	AV_PIX_FMT_BAYER_BGGR16BE               // bayer, BGBG..(odd line), GRGR..(even line), 16-bit samples, big-endian
	AV_PIX_FMT_BAYER_RGGB16LE               // bayer, RGRG..(odd line), GBGB..(even line), 16-bit samples, little-endian
	AV_PIX_FMT_BAYER_RGGB16BE               // bayer, RGRG..(odd line), GBGB..(even line), 16-bit samples, big-endian
	AV_PIX_FMT_BAYER_GBRG16LE               // bayer, GBGB..(odd line), RGRG..(even line), 16-bit samples, little-endian
	AV_PIX_FMT_BAYER_GBRG16BE               // bayer, GBGB..(odd line), RGRG..(even line), 16-bit samples, big-endian
	AV_PIX_FMT_BAYER_GRBG16LE               // bayer, GRGR..(odd line), BGBG..(even line), 16-bit samples, little-endian
	AV_PIX_FMT_BAYER_GRBG16BE               // bayer, GRGR..(odd line), BGBG..(even line), 16-bit samples, big-endian
	AV_PIX_FMT_XVMC                         // XVideo Motion Acceleration via common packet passing.
	AV_PIX_FMT_YUV440P10LE                  // planar YUV 4:4:0,20bpp, (1 Cr & Cb sample per 1x2 Y samples), little-endian
	AV_PIX_FMT_YUV440P10BE                  // planar YUV 4:4:0,20bpp, (1 Cr & Cb sample per 1x2 Y samples), big-endian
	AV_PIX_FMT_YUV440P12LE                  // planar YUV 4:4:0,24bpp, (1 Cr & Cb sample per 1x2 Y samples), little-endian
	AV_PIX_FMT_YUV440P12BE                  // planar YUV 4:4:0,24bpp, (1 Cr & Cb sample per 1x2 Y samples), big-endian
	AV_PIX_FMT_AYUV64LE                     // packed AYUV 4:4:4,64bpp (1 Cr & Cb sample per 1x1 Y & A samples), little-endian
	AV_PIX_FMT_AYUV64BE                     // packed AYUV 4:4:4,64bpp (1 Cr & Cb sample per 1x1 Y & A samples), big-endian
	AV_PIX_FMT_VIDEOTOOLBOX                 // hardware decoding through Videotoolbox
	AV_PIX_FMT_P010LE                       // like NV12, with 10bpp per component, data in the high bits, zeros in the low bits, little-endian
	AV_PIX_FMT_P010BE                       // like NV12, with 10bpp per component, data in the high bits, zeros in the low bits, big-endian
	AV_PIX_FMT_GBRAP12BE                    // planar GBR 4:4:4:4 48bpp, big-endian
	AV_PIX_FMT_GBRAP12LE                    // planar GBR 4:4:4:4 48bpp, little-endian
	AV_PIX_FMT_GBRAP10BE                    // planar GBR 4:4:4:4 40bpp, big-endian
	AV_PIX_FMT_GBRAP10LE                    // planar GBR 4:4:4:4 40bpp, little-endian
	AV_PIX_FMT_MEDIACODEC                   // hardware decoding through MediaCodec
	AV_PIX_FMT_GRAY12BE                     // Y , 12bpp, big-endian.
	AV_PIX_FMT_GRAY12LE                     // Y , 12bpp, little-endian.
	AV_PIX_FMT_GRAY10BE                     // Y , 10bpp, big-endian.
	AV_PIX_FMT_GRAY10LE                     // Y , 10bpp, little-endian.
	AV_PIX_FMT_P016LE                       // like NV12, with 16bpp per component, little-endian
	AV_PIX_FMT_P016BE                       // like NV12, with 16bpp per component, big-endian
	AV_PIX_FMT_D3D11                        // Hardware surfaces for Direct3D11.
	AV_PIX_FMT_GRAY9BE                      // Y , 9bpp, big-endian.
	AV_PIX_FMT_GRAY9LE                      // Y , 9bpp, little-endian.
	AV_PIX_FMT_GBRPF32BE                    // IEEE-754 single precision planar GBR 4:4:4, 96bpp, big-endian.
	AV_PIX_FMT_GBRPF32LE                    // IEEE-754 single precision planar GBR 4:4:4, 96bpp, little-endian.
	AV_PIX_FMT_GBRAPF32BE                   // IEEE-754 single precision planar GBRA 4:4:4:4, 128bpp, big-endian.
	AV_PIX_FMT_GBRAPF32LE                   // IEEE-754 single precision planar GBRA 4:4:4:4, 128bpp, little-endian.
	AV_PIX_FMT_DRM_PRIME                    // DRM-managed buffers exposed through PRIME buffer sharing.
	AV_PIX_FMT_OPENCL                       // Hardware surfaces for OpenCL.
	AV_PIX_FMT_GRAY14BE                     // Y , 14bpp, big-endian.
	AV_PIX_FMT_GRAY14LE                     // Y , 14bpp, little-endian.
	AV_PIX_FMT_GRAYF32BE                    // IEEE-754 single precision Y, 32bpp, big-endian.
	AV_PIX_FMT_GRAYF32LE                    // IEEE-754 single precision Y, 32bpp, little-endian.
)

const (
	AV_SAMPLE_FMT_NONE AVSampleFormat = iota
	AV_SAMPLE_FMT_U8                  //	unsigned 8 bits
	AV_SAMPLE_FMT_S16                 //	signed 16 bits
	AV_SAMPLE_FMT_S32                 //	signed 32 bits
	AV_SAMPLE_FMT_FLT                 //	float
	AV_SAMPLE_FMT_DBL                 //	double
	AV_SAMPLE_FMT_U8P                 //	unsigned 8 bits, planar
	AV_SAMPLE_FMT_S16P                //	signed 16 bits, planar
	AV_SAMPLE_FMT_S32P                //	signed 32 bits, planar
	AV_SAMPLE_FMT_FLTP                //	float, planar
	AV_SAMPLE_FMT_DBLP                //	double, planar
	AV_SAMPLE_FMT_S64                 //	signed 64 bits
	AV_SAMPLE_FMT_S64P                //	signed 64 bits, planar
)

func (f AVPixelFormat) String() string {
	switch f {
	case AV_PIX_FMT_NONE:
		return "AV_PIX_FMT_NONE"
	case AV_PIX_FMT_YUV420P:
		return "AV_PIX_FMT_YUV420P"
	case AV_PIX_FMT_YUYV422:
		return "AV_PIX_FMT_YUYV422"
	case AV_PIX_FMT_RGB24:
		return "AV_PIX_FMT_RGB24"
	case AV_PIX_FMT_BGR24:
		return "AV_PIX_FMT_BGR24"
	case AV_PIX_FMT_YUV422P:
		return "AV_PIX_FMT_YUV422P"
	case AV_PIX_FMT_YUV444P:
		return "AV_PIX_FMT_YUV444P"
	case AV_PIX_FMT_YUV410P:
		return "AV_PIX_FMT_YUV410P"
	case AV_PIX_FMT_YUV411P:
		return "AV_PIX_FMT_YUV411P"
	case AV_PIX_FMT_GRAY8:
		return "AV_PIX_FMT_GRAY8"
	case AV_PIX_FMT_MONOWHITE:
		return "AV_PIX_FMT_MONOWHITE"
	case AV_PIX_FMT_MONOBLACK:
		return "AV_PIX_FMT_MONOBLACK"
	case AV_PIX_FMT_PAL8:
		return "AV_PIX_FMT_PAL8"
	case AV_PIX_FMT_YUVJ420P:
		return "AV_PIX_FMT_YUVJ420P"
	case AV_PIX_FMT_YUVJ422P:
		return "AV_PIX_FMT_YUVJ422P"
	case AV_PIX_FMT_YUVJ444P:
		return "AV_PIX_FMT_YUVJ444P"
	case AV_PIX_FMT_UYVY422:
		return "AV_PIX_FMT_UYVY422"
	case AV_PIX_FMT_UYYVYY411:
		return "AV_PIX_FMT_UYYVYY411"
	case AV_PIX_FMT_BGR8:
		return "AV_PIX_FMT_BGR8"
	case AV_PIX_FMT_BGR4:
		return "AV_PIX_FMT_BGR4"
	case AV_PIX_FMT_BGR4_BYTE:
		return "AV_PIX_FMT_BGR4_BYTE"
	case AV_PIX_FMT_RGB8:
		return "AV_PIX_FMT_RGB8"
	case AV_PIX_FMT_RGB4:
		return "AV_PIX_FMT_RGB4"
	case AV_PIX_FMT_RGB4_BYTE:
		return "AV_PIX_FMT_RGB4_BYTE"
	case AV_PIX_FMT_NV12:
		return "AV_PIX_FMT_NV12"
	case AV_PIX_FMT_NV21:
		return "AV_PIX_FMT_NV21"
	case AV_PIX_FMT_ARGB:
		return "AV_PIX_FMT_ARGB"
	case AV_PIX_FMT_RGBA:
		return "AV_PIX_FMT_RGBA"
	case AV_PIX_FMT_ABGR:
		return "AV_PIX_FMT_ABGR"
	case AV_PIX_FMT_BGRA:
		return "AV_PIX_FMT_BGRA"
	case AV_PIX_FMT_GRAY16BE:
		return "AV_PIX_FMT_GRAY16BE"
	case AV_PIX_FMT_GRAY16LE:
		return "AV_PIX_FMT_GRAY16LE"
	case AV_PIX_FMT_YUV440P:
		return "AV_PIX_FMT_YUV440P"
	case AV_PIX_FMT_YUVJ440P:
		return "AV_PIX_FMT_YUVJ440P"
	case AV_PIX_FMT_YUVA420P:
		return "AV_PIX_FMT_YUVA420P"
	case AV_PIX_FMT_RGB48BE:
		return "AV_PIX_FMT_RGB48BE"
	case AV_PIX_FMT_RGB48LE:
		return "AV_PIX_FMT_RGB48LE"
	case AV_PIX_FMT_RGB565BE:
		return "AV_PIX_FMT_RGB565BE"
	case AV_PIX_FMT_RGB565LE:
		return "AV_PIX_FMT_RGB565LE"
	case AV_PIX_FMT_RGB555BE:
		return "AV_PIX_FMT_RGB555BE"
	case AV_PIX_FMT_RGB555LE:
		return "AV_PIX_FMT_RGB555LE"
	case AV_PIX_FMT_BGR565BE:
		return "AV_PIX_FMT_BGR565BE"
	case AV_PIX_FMT_BGR565LE:
		return "AV_PIX_FMT_BGR565LE"
	case AV_PIX_FMT_BGR555BE:
		return "AV_PIX_FMT_BGR555BE"
	case AV_PIX_FMT_BGR555LE:
		return "AV_PIX_FMT_BGR555LE"
	case AV_PIX_FMT_VAAPI_MOCO:
		return "AV_PIX_FMT_VAAPI_MOCO"
	case AV_PIX_FMT_VAAPI_IDCT:
		return "AV_PIX_FMT_VAAPI_IDCT"
	case AV_PIX_FMT_VAAPI_VLD:
		return "AV_PIX_FMT_VAAPI_VLD"
	case AV_PIX_FMT_VAAPI:
		return "AV_PIX_FMT_VAAPI"
	case AV_PIX_FMT_YUV420P16LE:
		return "AV_PIX_FMT_YUV420P16LE"
	case AV_PIX_FMT_YUV420P16BE:
		return "AV_PIX_FMT_YUV420P16BE"
	case AV_PIX_FMT_YUV422P16LE:
		return "AV_PIX_FMT_YUV422P16LE"
	case AV_PIX_FMT_YUV422P16BE:
		return "AV_PIX_FMT_YUV422P16BE"
	case AV_PIX_FMT_YUV444P16LE:
		return "AV_PIX_FMT_YUV444P16LE"
	case AV_PIX_FMT_YUV444P16BE:
		return "AV_PIX_FMT_YUV444P16BE"
	case AV_PIX_FMT_DXVA2_VLD:
		return "AV_PIX_FMT_DXVA2_VLD"
	case AV_PIX_FMT_RGB444LE:
		return "AV_PIX_FMT_RGB444LE"
	case AV_PIX_FMT_RGB444BE:
		return "AV_PIX_FMT_RGB444BE"
	case AV_PIX_FMT_BGR444LE:
		return "AV_PIX_FMT_BGR444LE"
	case AV_PIX_FMT_BGR444BE:
		return "AV_PIX_FMT_BGR444BE"
	case AV_PIX_FMT_YA8:
		return "AV_PIX_FMT_YA8"
	case AV_PIX_FMT_Y400A:
		return "AV_PIX_FMT_Y400A"
	case AV_PIX_FMT_GRAY8A:
		return "AV_PIX_FMT_GRAY8A"
	case AV_PIX_FMT_BGR48BE:
		return "AV_PIX_FMT_BGR48BE"
	case AV_PIX_FMT_BGR48LE:
		return "AV_PIX_FMT_BGR48LE"
	case AV_PIX_FMT_YUV420P9BE:
		return "AV_PIX_FMT_YUV420P9BE"
	case AV_PIX_FMT_YUV420P9LE:
		return "AV_PIX_FMT_YUV420P9LE"
	case AV_PIX_FMT_YUV420P10BE:
		return "AV_PIX_FMT_YUV420P10BE"
	case AV_PIX_FMT_YUV420P10LE:
		return "AV_PIX_FMT_YUV420P10LE"
	case AV_PIX_FMT_YUV422P10BE:
		return "AV_PIX_FMT_YUV422P10BE"
	case AV_PIX_FMT_YUV422P10LE:
		return "AV_PIX_FMT_YUV422P10LE"
	case AV_PIX_FMT_YUV444P9BE:
		return "AV_PIX_FMT_YUV444P9BE"
	case AV_PIX_FMT_YUV444P9LE:
		return "AV_PIX_FMT_YUV444P9LE"
	case AV_PIX_FMT_YUV444P10BE:
		return "AV_PIX_FMT_YUV444P10BE"
	case AV_PIX_FMT_YUV444P10LE:
		return "AV_PIX_FMT_YUV444P10LE"
	case AV_PIX_FMT_YUV422P9BE:
		return "AV_PIX_FMT_YUV422P9BE"
	case AV_PIX_FMT_YUV422P9LE:
		return "AV_PIX_FMT_YUV422P9LE"
	case AV_PIX_FMT_GBRP:
		return "AV_PIX_FMT_GBRP"
	case AV_PIX_FMT_GBR24P:
		return "AV_PIX_FMT_GBR24P"
	case AV_PIX_FMT_GBRP9BE:
		return "AV_PIX_FMT_GBRP9BE"
	case AV_PIX_FMT_GBRP9LE:
		return "AV_PIX_FMT_GBRP9LE"
	case AV_PIX_FMT_GBRP10BE:
		return "AV_PIX_FMT_GBRP10BE"
	case AV_PIX_FMT_GBRP10LE:
		return "AV_PIX_FMT_GBRP10LE"
	case AV_PIX_FMT_GBRP16BE:
		return "AV_PIX_FMT_GBRP16BE"
	case AV_PIX_FMT_GBRP16LE:
		return "AV_PIX_FMT_GBRP16LE"
	case AV_PIX_FMT_YUVA422P:
		return "AV_PIX_FMT_YUVA422P"
	case AV_PIX_FMT_YUVA444P:
		return "AV_PIX_FMT_YUVA444P"
	case AV_PIX_FMT_YUVA420P9BE:
		return "AV_PIX_FMT_YUVA420P9BE"
	case AV_PIX_FMT_YUVA420P9LE:
		return "AV_PIX_FMT_YUVA420P9LE"
	case AV_PIX_FMT_YUVA422P9BE:
		return "AV_PIX_FMT_YUVA422P9BE"
	case AV_PIX_FMT_YUVA422P9LE:
		return "AV_PIX_FMT_YUVA422P9LE"
	case AV_PIX_FMT_YUVA444P9BE:
		return "AV_PIX_FMT_YUVA444P9BE"
	case AV_PIX_FMT_YUVA444P9LE:
		return "AV_PIX_FMT_YUVA444P9LE"
	case AV_PIX_FMT_YUVA420P10BE:
		return "AV_PIX_FMT_YUVA420P10BE"
	case AV_PIX_FMT_YUVA420P10LE:
		return "AV_PIX_FMT_YUVA420P10LE"
	case AV_PIX_FMT_YUVA422P10BE:
		return "AV_PIX_FMT_YUVA422P10BE"
	case AV_PIX_FMT_YUVA422P10LE:
		return "AV_PIX_FMT_YUVA422P10LE"
	case AV_PIX_FMT_YUVA444P10BE:
		return "AV_PIX_FMT_YUVA444P10BE"
	case AV_PIX_FMT_YUVA444P10LE:
		return "AV_PIX_FMT_YUVA444P10LE"
	case AV_PIX_FMT_YUVA420P16BE:
		return "AV_PIX_FMT_YUVA420P16BE"
	case AV_PIX_FMT_YUVA420P16LE:
		return "AV_PIX_FMT_YUVA420P16LE"
	case AV_PIX_FMT_YUVA422P16BE:
		return "AV_PIX_FMT_YUVA422P16BE"
	case AV_PIX_FMT_YUVA422P16LE:
		return "AV_PIX_FMT_YUVA422P16LE"
	case AV_PIX_FMT_YUVA444P16BE:
		return "AV_PIX_FMT_YUVA444P16BE"
	case AV_PIX_FMT_YUVA444P16LE:
		return "AV_PIX_FMT_YUVA444P16LE"
	case AV_PIX_FMT_VDPAU:
		return "AV_PIX_FMT_VDPAU"
	case AV_PIX_FMT_XYZ12LE:
		return "AV_PIX_FMT_XYZ12LE"
	case AV_PIX_FMT_XYZ12BE:
		return "AV_PIX_FMT_XYZ12BE"
	case AV_PIX_FMT_NV16:
		return "AV_PIX_FMT_NV16"
	case AV_PIX_FMT_NV20LE:
		return "AV_PIX_FMT_NV20LE"
	case AV_PIX_FMT_NV20BE:
		return "AV_PIX_FMT_NV20BE"
	case AV_PIX_FMT_RGBA64BE:
		return "AV_PIX_FMT_RGBA64BE"
	case AV_PIX_FMT_RGBA64LE:
		return "AV_PIX_FMT_RGBA64LE"
	case AV_PIX_FMT_BGRA64BE:
		return "AV_PIX_FMT_BGRA64BE"
	case AV_PIX_FMT_BGRA64LE:
		return "AV_PIX_FMT_BGRA64LE"
	case AV_PIX_FMT_YVYU422:
		return "AV_PIX_FMT_YVYU422"
	case AV_PIX_FMT_YA16BE:
		return "AV_PIX_FMT_YA16BE"
	case AV_PIX_FMT_YA16LE:
		return "AV_PIX_FMT_YA16LE"
	case AV_PIX_FMT_GBRAP:
		return "AV_PIX_FMT_GBRAP"
	case AV_PIX_FMT_GBRAP16BE:
		return "AV_PIX_FMT_GBRAP16BE"
	case AV_PIX_FMT_GBRAP16LE:
		return "AV_PIX_FMT_GBRAP16LE"
	case AV_PIX_FMT_QSV:
		return "AV_PIX_FMT_QSV"
	case AV_PIX_FMT_MMAL:
		return "AV_PIX_FMT_MMAL"
	case AV_PIX_FMT_D3D11VA_VLD:
		return "AV_PIX_FMT_D3D11VA_VLD"
	case AV_PIX_FMT_CUDA:
		return "AV_PIX_FMT_CUDA"
	case AV_PIX_FMT_0RGB:
		return "AV_PIX_FMT_0RGB"
	case AV_PIX_FMT_RGB0:
		return "AV_PIX_FMT_RGB0"
	case AV_PIX_FMT_0BGR:
		return "AV_PIX_FMT_0BGR"
	case AV_PIX_FMT_BGR0:
		return "AV_PIX_FMT_BGR0"
	case AV_PIX_FMT_YUV420P12BE:
		return "AV_PIX_FMT_YUV420P12BE"
	case AV_PIX_FMT_YUV420P12LE:
		return "AV_PIX_FMT_YUV420P12LE"
	case AV_PIX_FMT_YUV420P14BE:
		return "AV_PIX_FMT_YUV420P14BE"
	case AV_PIX_FMT_YUV420P14LE:
		return "AV_PIX_FMT_YUV420P14LE"
	case AV_PIX_FMT_YUV422P12BE:
		return "AV_PIX_FMT_YUV422P12BE"
	case AV_PIX_FMT_YUV422P12LE:
		return "AV_PIX_FMT_YUV422P12LE"
	case AV_PIX_FMT_YUV422P14BE:
		return "AV_PIX_FMT_YUV422P14BE"
	case AV_PIX_FMT_YUV422P14LE:
		return "AV_PIX_FMT_YUV422P14LE"
	case AV_PIX_FMT_YUV444P12BE:
		return "AV_PIX_FMT_YUV444P12BE"
	case AV_PIX_FMT_YUV444P12LE:
		return "AV_PIX_FMT_YUV444P12LE"
	case AV_PIX_FMT_YUV444P14BE:
		return "AV_PIX_FMT_YUV444P14BE"
	case AV_PIX_FMT_YUV444P14LE:
		return "AV_PIX_FMT_YUV444P14LE"
	case AV_PIX_FMT_GBRP12BE:
		return "AV_PIX_FMT_GBRP12BE"
	case AV_PIX_FMT_GBRP12LE:
		return "AV_PIX_FMT_GBRP12LE"
	case AV_PIX_FMT_GBRP14BE:
		return "AV_PIX_FMT_GBRP14BE"
	case AV_PIX_FMT_GBRP14LE:
		return "AV_PIX_FMT_GBRP14LE"
	case AV_PIX_FMT_YUVJ411P:
		return "AV_PIX_FMT_YUVJ411P"
	case AV_PIX_FMT_BAYER_BGGR8:
		return "AV_PIX_FMT_BAYER_BGGR8"
	case AV_PIX_FMT_BAYER_RGGB8:
		return "AV_PIX_FMT_BAYER_RGGB8"
	case AV_PIX_FMT_BAYER_GBRG8:
		return "AV_PIX_FMT_BAYER_GBRG8"
	case AV_PIX_FMT_BAYER_GRBG8:
		return "AV_PIX_FMT_BAYER_GRBG8"
	case AV_PIX_FMT_BAYER_BGGR16LE:
		return "AV_PIX_FMT_BAYER_BGGR16LE"
	case AV_PIX_FMT_BAYER_BGGR16BE:
		return "AV_PIX_FMT_BAYER_BGGR16BE"
	case AV_PIX_FMT_BAYER_RGGB16LE:
		return "AV_PIX_FMT_BAYER_RGGB16LE"
	case AV_PIX_FMT_BAYER_RGGB16BE:
		return "AV_PIX_FMT_BAYER_RGGB16BE"
	case AV_PIX_FMT_BAYER_GBRG16LE:
		return "AV_PIX_FMT_BAYER_GBRG16LE"
	case AV_PIX_FMT_BAYER_GBRG16BE:
		return "AV_PIX_FMT_BAYER_GBRG16BE"
	case AV_PIX_FMT_BAYER_GRBG16LE:
		return "AV_PIX_FMT_BAYER_GRBG16LE"
	case AV_PIX_FMT_BAYER_GRBG16BE:
		return "AV_PIX_FMT_BAYER_GRBG16BE"
	case AV_PIX_FMT_XVMC:
		return "AV_PIX_FMT_XVMC"
	case AV_PIX_FMT_YUV440P10LE:
		return "AV_PIX_FMT_YUV440P10LE"
	case AV_PIX_FMT_YUV440P10BE:
		return "AV_PIX_FMT_YUV440P10BE"
	case AV_PIX_FMT_YUV440P12LE:
		return "AV_PIX_FMT_YUV440P12LE"
	case AV_PIX_FMT_YUV440P12BE:
		return "AV_PIX_FMT_YUV440P12BE"
	case AV_PIX_FMT_AYUV64LE:
		return "AV_PIX_FMT_AYUV64LE"
	case AV_PIX_FMT_AYUV64BE:
		return "AV_PIX_FMT_AYUV64BE"
	case AV_PIX_FMT_VIDEOTOOLBOX:
		return "AV_PIX_FMT_VIDEOTOOLBOX"
	case AV_PIX_FMT_P010LE:
		return "AV_PIX_FMT_P010LE"
	case AV_PIX_FMT_P010BE:
		return "AV_PIX_FMT_P010BE"
	case AV_PIX_FMT_GBRAP12BE:
		return "AV_PIX_FMT_GBRAP12BE"
	case AV_PIX_FMT_GBRAP12LE:
		return "AV_PIX_FMT_GBRAP12LE"
	case AV_PIX_FMT_GBRAP10BE:
		return "AV_PIX_FMT_GBRAP10BE"
	case AV_PIX_FMT_GBRAP10LE:
		return "AV_PIX_FMT_GBRAP10LE"
	case AV_PIX_FMT_MEDIACODEC:
		return "AV_PIX_FMT_MEDIACODEC"
	case AV_PIX_FMT_GRAY12BE:
		return "AV_PIX_FMT_GRAY12BE"
	case AV_PIX_FMT_GRAY12LE:
		return "AV_PIX_FMT_GRAY12LE"
	case AV_PIX_FMT_GRAY10BE:
		return "AV_PIX_FMT_GRAY10BE"
	case AV_PIX_FMT_GRAY10LE:
		return "AV_PIX_FMT_GRAY10LE"
	case AV_PIX_FMT_P016LE:
		return "AV_PIX_FMT_P016LE"
	case AV_PIX_FMT_P016BE:
		return "AV_PIX_FMT_P016BE"
	case AV_PIX_FMT_D3D11:
		return "AV_PIX_FMT_D3D11"
	case AV_PIX_FMT_GRAY9BE:
		return "AV_PIX_FMT_GRAY9BE"
	case AV_PIX_FMT_GRAY9LE:
		return "AV_PIX_FMT_GRAY9LE"
	case AV_PIX_FMT_GBRPF32BE:
		return "AV_PIX_FMT_GBRPF32BE"
	case AV_PIX_FMT_GBRPF32LE:
		return "AV_PIX_FMT_GBRPF32LE"
	case AV_PIX_FMT_GBRAPF32BE:
		return "AV_PIX_FMT_GBRAPF32BE"
	case AV_PIX_FMT_GBRAPF32LE:
		return "AV_PIX_FMT_GBRAPF32LE"
	case AV_PIX_FMT_DRM_PRIME:
		return "AV_PIX_FMT_DRM_PRIME"
	case AV_PIX_FMT_OPENCL:
		return "AV_PIX_FMT_OPENCL"
	case AV_PIX_FMT_GRAY14BE:
		return "AV_PIX_FMT_GRAY14BE"
	case AV_PIX_FMT_GRAY14LE:
		return "AV_PIX_FMT_GRAY14LE"
	case AV_PIX_FMT_GRAYF32BE:
		return "AV_PIX_FMT_GRAYF32BE"
	case AV_PIX_FMT_GRAYF32LE:
		return "AV_PIX_FMT_GRAYF32LE"
	default:
		return "[?? Invalid AVPixelFormat value]"
	}
}

func (f AVSampleFormat) String() string {
	switch f {
	case AV_SAMPLE_FMT_NONE:
		return "AV_SAMPLE_FMT_NONE"
	case AV_SAMPLE_FMT_U8:
		return "AV_SAMPLE_FMT_U8"
	case AV_SAMPLE_FMT_S16:
		return "AV_SAMPLE_FMT_S16"
	case AV_SAMPLE_FMT_S32:
		return "AV_SAMPLE_FMT_S32"
	case AV_SAMPLE_FMT_FLT:
		return "AV_SAMPLE_FMT_FLT"
	case AV_SAMPLE_FMT_DBL:
		return "AV_SAMPLE_FMT_DBL"
	case AV_SAMPLE_FMT_U8P:
		return "AV_SAMPLE_FMT_U8P"
	case AV_SAMPLE_FMT_S16P:
		return "AV_SAMPLE_FMT_S16P"
	case AV_SAMPLE_FMT_S32P:
		return "AV_SAMPLE_FMT_S32P"
	case AV_SAMPLE_FMT_FLTP:
		return "AV_SAMPLE_FMT_FLTP"
	case AV_SAMPLE_FMT_DBLP:
		return "AV_SAMPLE_FMT_DBLP"
	case AV_SAMPLE_FMT_S64:
		return "AV_SAMPLE_FMT_S64"
	case AV_SAMPLE_FMT_S64P:
		return "AV_SAMPLE_FMT_S64P"
	default:
		return "[?? Invalid AVSampleFormat value]"
	}
}
