package waveshare

import (
	"context"
	"fmt"
	"image"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/gpio/broadcom"
)

type EPD struct {
	gopi.Unit
	gopi.GPIO
	gopi.SPI
}

const (
	EPD_PIN_RESET = gopi.GPIOPin(17)
	EPD_PIN_CS    = gopi.GPIOPin(8)
	EPD_PIN_DC    = gopi.GPIOPin(25)
	EPD_PIN_BUSY  = gopi.GPIOPin(24)

	EPD_SPI_SPEED = 10000000
	EPD_SPI_BUS   = 0
	EPD_SPI_SLAVE = 0
	EPD_SPI_MODE  = gopi.SPI_MODE_0
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *EPD) Define(cfg gopi.Config) {
	cfg.FlagUint("epd.width", 880, "Width of display")
	cfg.FlagUint("epd.height", 528, "Height of display")
}

func (this *EPD) New(gopi.Config) error {
	if this.GPIO == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing GPIO interface")
	}
	if this.SPI == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing SPI interface")
	}
	return nil
}

func (this *EPD) Dispose() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *EPD) Init() error {
	//this.reset()
	/*
	   EPD_7IN5_HD_WaitUntilIdle();
	   EPD_7IN5_HD_SendCommand(0x12);  //SWRESET
	   EPD_7IN5_HD_WaitUntilIdle();

	   EPD_7IN5_HD_SendCommand(0x46);  // Auto Write Red RAM
	   EPD_7IN5_HD_SendData(0xf7);
	   EPD_7IN5_HD_WaitUntilIdle();
	   EPD_7IN5_HD_SendCommand(0x47);  // Auto Write  B/W RAM
	   EPD_7IN5_HD_SendData(0xf7);
	   EPD_7IN5_HD_WaitUntilIdle();


	   EPD_7IN5_HD_SendCommand(0x0C);  // Soft start setting
	   EPD_7IN5_HD_SendData(0xAE);
	   EPD_7IN5_HD_SendData(0xC7);
	   EPD_7IN5_HD_SendData(0xC3);
	   EPD_7IN5_HD_SendData(0xC0);
	   EPD_7IN5_HD_SendData(0x40);


	   EPD_7IN5_HD_SendCommand(0x01);  // Set MUX as 527
	   EPD_7IN5_HD_SendData(0xAF);
	   EPD_7IN5_HD_SendData(0x02);
	   EPD_7IN5_HD_SendData(0x01);//0x01


	   EPD_7IN5_HD_SendCommand(0x11);  // Data entry mode
	   EPD_7IN5_HD_SendData(0x01);

	   EPD_7IN5_HD_SendCommand(0x44);
	   EPD_7IN5_HD_SendData(0x00); // RAM x address start at 0
	   EPD_7IN5_HD_SendData(0x00);
	   EPD_7IN5_HD_SendData(0x6F);
	   EPD_7IN5_HD_SendData(0x03);
	   EPD_7IN5_HD_SendCommand(0x45);
	   EPD_7IN5_HD_SendData(0xFF);
	   EPD_7IN5_HD_SendData(0x03);
	   EPD_7IN5_HD_SendData(0x00);
	   EPD_7IN5_HD_SendData(0x00);


	   EPD_7IN5_HD_SendCommand(0x3C); // VBD
	   EPD_7IN5_HD_SendData(0x05); // LUT1, for white

	   EPD_7IN5_HD_SendCommand(0x18);
	   EPD_7IN5_HD_SendData(0X80);


	   EPD_7IN5_HD_SendCommand(0x22);
	   EPD_7IN5_HD_SendData(0XB1); //Load Temperature and waveform setting.
	   EPD_7IN5_HD_SendCommand(0x20);
	   EPD_7IN5_HD_WaitUntilIdle();

	   EPD_7IN5_HD_SendCommand(0x4E); // set RAM x address count to 0;
	   EPD_7IN5_HD_SendData(0x00);
	   EPD_7IN5_HD_SendData(0x00);
	   EPD_7IN5_HD_SendCommand(0x4F);
	   EPD_7IN5_HD_SendData(0x00);
	   EPD_7IN5_HD_SendData(0x00);
	*/

	// Return sucess
	return nil
}

func (this *EPD) Clear() error {
	/*
	       UDOUBLE Width, Height;
	       Width =(EPD_7IN5_HD_WIDTH % 8 == 0)?(EPD_7IN5_HD_WIDTH / 8 ):(EPD_7IN5_HD_WIDTH / 8 + 1);
	       Height = EPD_7IN5_HD_HEIGHT;

	       EPD_7IN5_HD_SendCommand(0x4F);
	       EPD_7IN5_HD_SendData(0x00);
	       EPD_7IN5_HD_SendData(0x00);
	        EPD_7IN5_HD_SendCommand(0x24);
	       UDOUBLE i;
	       for(i=0; i<58080; i++) {
	           EPD_7IN5_HD_SendData(0xff);
	       }

	       EPD_7IN5_HD_SendCommand(0x26);
	       for(i=0; i<Height*Width; i++){
	           EPD_7IN5_HD_SendData(0xff);
	       }
	       EPD_7IN5_HD_SendCommand(0x22);
	       EPD_7IN5_HD_SendData(0xF7);//Load LUT from MCU(0x32)
	       EPD_7IN5_HD_SendCommand(0x20);
	       DEV_Delay_ms(10);


	   	EPD_7IN5_HD_WaitUntilIdle();
	*/
	return nil
}

func (this *EPD) Display(image.Image) error {
	/*
	   ******************************************************************************
	   function :	Sends the image buffer in RAM to e-Paper and displays
	   parameter:
	   ******************************************************************************
	   void EPD_7IN5_HD_Display(const UBYTE *blackimage)
	   {
	       UDOUBLE Width, Height;
	       Width =(EPD_7IN5_HD_WIDTH % 8 == 0)?(EPD_7IN5_HD_WIDTH / 8 ):(EPD_7IN5_HD_WIDTH / 8 + 1);
	       Height = EPD_7IN5_HD_HEIGHT;

	       EPD_7IN5_HD_SendCommand(0x4F);
	       EPD_7IN5_HD_SendData(0x00);
	       EPD_7IN5_HD_SendData(0x00);
	       EPD_7IN5_HD_SendCommand(0x24);

	       UDOUBLE i;
	       for (UDOUBLE j = 0; j < Height; j++) {
	           for (UDOUBLE i = 0; i < Width; i++) {
	               EPD_7IN5_HD_SendData(blackimage[i + j * Width]);
	           }
	       }

	       EPD_7IN5_HD_SendCommand(0x26);
	       for(i=0; i<Height*Width; i++)	{
	           EPD_7IN5_HD_SendData(0xff);
	       }
	       EPD_7IN5_HD_SendCommand(0x22);
	       EPD_7IN5_HD_SendData(0xF7);//Load LUT from MCU(0x32)
	       EPD_7IN5_HD_SendCommand(0x20);
	       DEV_Delay_ms(10);
	       EPD_7IN5_HD_WaitUntilIdle();
	   }
	*/
	// Return success
	return nil
}

func (this *EPD) Sleep() {
	this.sendCommand(0x10)
	this.sendData(0x01)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// waitUntilIdle waits until busy pin goes low
func (this *EPD) waitUntilIdle(ctx context.Context) error {
	ticker := time.NewTimer(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if this.GPIO.ReadPin(EPD_PIN_BUSY) == gopi.GPIO_LOW {
				// DEV_Delay_ms(200);
				return nil
			}
			ticker.Reset(10 * time.Millisecond)
		}
	}
}

func (this *EPD) sendCommand(reg uint8) {
	/*
		DEV_Digital_Write(EPD_DC_PIN, 0)
		DEV_Digital_Write(EPD_CS_PIN, 0)
		DEV_SPI_WriteByte(Reg)
		DEV_Digital_Write(EPD_CS_PIN, 1)*/
}

func (this *EPD) sendData(reg uint8) {
	/*
	   DEV_Digital_Write(EPD_DC_PIN, 1);
	   DEV_Digital_Write(EPD_CS_PIN, 0);
	   DEV_SPI_WriteByte(Data);
	   DEV_Digital_Write(EPD_CS_PIN, 1);*/
}

/*

static void EPD_7IN5_HD_Reset(void)
{
    DEV_Digital_Write(EPD_RST_PIN, 1);
    DEV_Delay_ms(200);
    DEV_Digital_Write(EPD_RST_PIN, 0);
    DEV_Delay_ms(2);
    DEV_Digital_Write(EPD_RST_PIN, 1);
    DEV_Delay_ms(200);
}

function : send command
parameter:
     Reg : Command register

function :	send data
parameter:
    Data : Write data
static void EPD_7IN5_HD_S
function :	Wait until the busy_pin goes LOW
parameter:
static void EPD_7IN5_HD_WaitUntilIdle(void)
{
    Debug("e-Paper busy\r\n");
    do{
        DEV_Delay_ms(10);
    }while(DEV_Digital_Read(EPD_BUSY_PIN) == 1);
    DEV_Delay_ms(200);
    Debug("e-Paper busy release\r\n");

}
*/

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *EPD) String() string {
	str := "<epd"
	str += " gpio=" + fmt.Sprint(this.GPIO)
	return str + ">"
}
