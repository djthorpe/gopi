/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// BME280
//
// This package implements the interface to the temperature, pressure
// and humidity sensor BME280 over the I2C interface. You should start by
// creating an instance to communicate with the sensor:
//
//    bme280, err := gopi.Open(adafruit.BME280{ /* configuration */ })
//    if err != nil { /* handle errors */ }
//    defer bme280.Close()
//
// You can then read temperature, atmospheric pressure and humidity values:
//
//    temp, atmospheric, humidity, err := bme280.(*adafruit.BME280Driver).ReadValues()
//    if err != nil { /* handle errors */ }
//
// The temperature is provided in Celcius, the pressure in hectopascal hPa and
// humidity in percentage of relative humidity, %RH. You can also determine the
// altitude above sea level, if you provide the pressure value at sealevel:
//
//    altitude := bme280.(*BME280Driver).AltitudeForPressure(atmospheric,sealevel)
//
// You can use the constant adafruit.BME280_PRESSURE_SEALEVEL for the sealevel
// pressure value.
//
// By default, the configuration only needs an I2C parameter, but you can also
// set the Slave parameter if you are able to alter the slave address of the
// device on the I2C bus:
//
//   config := adafruit.BME280{ I2C: i2c, Slave: addr }
//   bme280, err := gopi.Open(config)
//   if err == adafruit.ErrNoDevice { /* No BME280 device detected */ }
//
package adafruit

import (
	"errors"
	"fmt"
	"math"
)

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

type BME280 struct {
	// The I2C Hardware
	I2C hw.I2CDriver

	// The slave address, usually 0x77 or 0x76
	Slave uint8
}

type BME280Driver struct {
	i2c         hw.I2CDriver
	slave       uint8
	chipid      uint8
	calibration *BME280Calibation
	oversample  BME280Oversample
	log         *util.LoggerDevice
}

// The calibation values which are used to calculate
type BME280Calibation struct {
	T1                             uint16
	T2, T3                         int16
	P1                             uint16
	P2, P3, P4, P5, P6, P7, P8, P9 int16
	H1                             uint8
	H2                             int16
	H3                             uint8
	H4, H5                         int16
	H6                             int8
}

// Oversampling type
type BME280Oversample uint8

////////////////////////////////////////////////////////////////////////////////

const (
	BME280_SLAVE_DEFAULT uint8 = 0x77
	BME280_SOFTRESET_VALUE uint8 = 0xB6
	BME280_PRESSURE_SEALEVEL float64 = 1013.25
)

const (
	// Oversampling mode
	BME280_OVERSAMPLE_NONE BME280Oversample = iota
	BME280_OVERSAMPLE_1
	BME280_OVERSAMPLE_2
	BME280_OVERSAMPLE_4
	BME280_OVERSAMPLE_8
	BME280_OVERSAMPLE_16
)

const (
	BME280_REGISTER_DIG_T1       uint8 = 0x88
	BME280_REGISTER_DIG_T2       uint8 = 0x8A
	BME280_REGISTER_DIG_T3       uint8 = 0x8C
	BME280_REGISTER_DIG_P1       uint8 = 0x8E
	BME280_REGISTER_DIG_P2       uint8 = 0x90
	BME280_REGISTER_DIG_P3       uint8 = 0x92
	BME280_REGISTER_DIG_P4       uint8 = 0x94
	BME280_REGISTER_DIG_P5       uint8 = 0x96
	BME280_REGISTER_DIG_P6       uint8 = 0x98
	BME280_REGISTER_DIG_P7       uint8 = 0x9A
	BME280_REGISTER_DIG_P8       uint8 = 0x9C
	BME280_REGISTER_DIG_P9       uint8 = 0x9E
	BME280_REGISTER_DIG_H1       uint8 = 0xA1
	BME280_REGISTER_DIG_H2       uint8 = 0xE1
	BME280_REGISTER_DIG_H3       uint8 = 0xE3
	BME280_REGISTER_DIG_H4       uint8 = 0xE4
	BME280_REGISTER_DIG_H5       uint8 = 0xE5
	BME280_REGISTER_DIG_H6       uint8 = 0xE7
	BME280_REGISTER_CHIPID       uint8 = 0xD0
	BME280_REGISTER_VERSION      uint8 = 0xD1
	BME280_REGISTER_SOFTRESET    uint8 = 0xE0
	BME280_REGISTER_CAL26        uint8 = 0xE1 // R calibration stored in 0xE1-0xF0
	BME280_REGISTER_CONTROLHUMID uint8 = 0xF2
	BME280_REGISTER_CONTROL      uint8 = 0xF4
	BME280_REGISTER_CONFIG       uint8 = 0xF5
	BME280_REGISTER_PRESSUREDATA uint8 = 0xF7
	BME280_REGISTER_TEMPDATA     uint8 = 0xFA
	BME280_REGISTER_HUMIDDATA    uint8 = 0xFD
)

////////////////////////////////////////////////////////////////////////////////

var (
	ErrNoDevice = errors.New("No BME280 device detected")
)

////////////////////////////////////////////////////////////////////////////////

func (config BME280) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug2("<adafruit.BME280>Open")

	this := new(BME280Driver)
	if config.Slave == 0 {
		this.slave = BME280_SLAVE_DEFAULT
	} else {
		this.slave = config.Slave
	}
	this.i2c = config.I2C
	this.log = log

	// Detect slave
	detected, err := this.i2c.DetectSlave(this.slave)
	if err != nil {
		return nil, err
	}
	if detected == false {
		return nil, ErrNoDevice
	}

	// Set slave
	if err := this.i2c.SetSlave(this.slave); err != nil {
		return nil, err
	}

	// Read Chip ID register
	chipid, err := this.i2c.ReadUint8(BME280_REGISTER_CHIPID)
	if err != nil {
		return nil, err
	}
	this.chipid = chipid

	// Load calibration values
	this.calibration, err = this.readCalibration()
	if err != nil {
		return nil, err
	}

	// Set Mode

	// 16x oversampling
	if err := this.i2c.WriteUint8(BME280_REGISTER_CONTROLHUMID, 0x05); err != nil {
		return nil, err
	}

	// 16x oversampling, normal mode
	if err := this.i2c.WriteUint8(BME280_REGISTER_CONTROL, 0xB7); err != nil {
		return nil, err
	}

	return this, nil
}

func (this *BME280Driver) Close() error {
	this.log.Debug2("<adafruit.BME280>Close")
	// No resources need to be freed
	return nil
}

func (this *BME280Driver) readCalibration() (*BME280Calibation, error) {
	var err error

	calibration := new(BME280Calibation)

	// Read temperature calibration values
	if calibration.T1, err = this.i2c.ReadUint16(BME280_REGISTER_DIG_T1); err != nil {
		return nil, err
	}
	if calibration.T2, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_T2); err != nil {
		return nil, err
	}
	if calibration.T3, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_T3); err != nil {
		return nil, err
	}

	// Read pressure calibration values
	if calibration.P1, err = this.i2c.ReadUint16(BME280_REGISTER_DIG_P1); err != nil {
		return nil, err
	}
	if calibration.P2, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P2); err != nil {
		return nil, err
	}
	if calibration.P3, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P3); err != nil {
		return nil, err
	}
	if calibration.P4, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P4); err != nil {
		return nil, err
	}
	if calibration.P5, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P5); err != nil {
		return nil, err
	}
	if calibration.P6, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P6); err != nil {
		return nil, err
	}
	if calibration.P7, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P7); err != nil {
		return nil, err
	}
	if calibration.P8, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P8); err != nil {
		return nil, err
	}
	if calibration.P9, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_P9); err != nil {
		return nil, err
	}

	// Read humidity calibration values. H4 and H5 are treated slightly differently
	if calibration.H1, err = this.i2c.ReadUint8(BME280_REGISTER_DIG_H1); err != nil {
		return nil, err
	}
	if calibration.H2, err = this.i2c.ReadInt16(BME280_REGISTER_DIG_H2); err != nil {
		return nil, err
	}
	if calibration.H3, err = this.i2c.ReadUint8(BME280_REGISTER_DIG_H3); err != nil {
		return nil, err
	}
	h41, err := this.i2c.ReadUint8(BME280_REGISTER_DIG_H4)
	if err != nil {
		return nil, err
	}
	h42, err := this.i2c.ReadUint8(BME280_REGISTER_DIG_H4 + 1)
	if err != nil {
		return nil, err
	}
	h51, err := this.i2c.ReadUint8(BME280_REGISTER_DIG_H5)
	if err != nil {
		return nil, err
	}
	h52, err := this.i2c.ReadUint8(BME280_REGISTER_DIG_H5 + 1)
	if err != nil {
		return nil, err
	}

	calibration.H4 = (int16(h41) << 4) | (int16(h42) & 0x0F)
	calibration.H5 = ((int16(h51) & 0xF0) >> 4) | int16(h52 << 4)

	if calibration.H6, err = this.i2c.ReadInt8(BME280_REGISTER_DIG_H6); err != nil {
		return nil, err
	}

	// Return calibration values
	return calibration, nil
}

// Device is reset using the complete power-on-reset procedure
func (this *BME280Driver) SoftReset() error {
	return this.i2c.WriteUint8(BME280_REGISTER_SOFTRESET,BME280_SOFTRESET_VALUE)
}

// Set oversampling value
func (this *BME280Driver) SetOversampleMode(value BME280Oversample) {
	this.oversample = value
}

////////////////////////////////////////////////////////////////////////////////

// Read compensated values and return them. The units for each are
// C, hPa and %RH. Error is returned if the data could not
// be returned
func (this *BME280Driver) ReadValues() (float64, float64, float64, error) {
	temp, t_fine, err := this.readTemperature()
	if err != nil {
		return 0,0,0,err
	}
	pressure, err := this.readPressure(t_fine)
	if err != nil {
		return 0,0,0,err
	}
	humidity, err := this.readHumidity(t_fine)
	if err != nil {
		return 0,0,0,err
	}
	return temp,pressure,humidity,nil
}

// Returns altitude in metres based on pressure reading in Pascals, given
// the sealevel pressure in Pascals. You can use a standard value of
// BME280_PRESSURE_SEALEVEL for sealevel
func (this *BME280Driver) AltitudeForPressure(atmospheric,sealevel float64) float64 {
	return 44330.0 * (1.0 - math.Pow(atmospheric / sealevel, (1.0/5.255)))
}

////////////////////////////////////////////////////////////////////////////////

func (this *BME280Driver) String() string {
	return fmt.Sprintf("<adafruit.BME280>{ slave=%02X chipid=%02X calibration=%v oversample=%v }", this.slave, this.chipid, this.calibration, this.oversample)
}

func (this *BME280Calibation) String() string {
	return fmt.Sprintf("<adafruit.BME280Calibation>{ T1=%v T2=%v T3=%v P1=%v P2=%v P3=%v P4=%v P5=%v P6=%v P7=%v P8=%v P9=%v H1=%v H2=%v H3=%v H4=%v H5=%v H6=%v }", this.T1, this.T2, this.T3, this.P1, this.P2, this.P3, this.P4, this.P5, this.P6, this.P7, this.P8, this.P9, this.H1, this.H2, this.H3, this.H4, this.H5, this.H6 )
}

func (o BME280Oversample) String() string {
	switch(o) {
	case BME280_OVERSAMPLE_NONE:
		return "BME280_OVERSAMPLE_NONE"
	case BME280_OVERSAMPLE_1:
		return "BME280_OVERSAMPLE_1"
	case BME280_OVERSAMPLE_2:
		return "BME280_OVERSAMPLE_2"
	case BME280_OVERSAMPLE_4:
		return "BME280_OVERSAMPLE_4"
	case BME280_OVERSAMPLE_8:
		return "BME280_OVERSAMPLE_8"
	case BME280_OVERSAMPLE_16:
		return "BME280_OVERSAMPLE_16"
	default:
		return "[?? Invalid BME280Oversampling value]"
	}
}

////////////////////////////////////////////////////////////////////////////////

// Return compensated temperature in Celcius, and the t_fine value
func (this *BME280Driver) readTemperature() (float64, float64, error) {
	adc, err := this.readTemperatureRaw()
	if err != nil {
		return 0,0,err
	}
	var1 := (float64(adc) / 16384.0 - float64(this.calibration.T1) / 1024.0) * float64(this.calibration.T2)
	var2 := ((float64(adc) / 131072.0 - float64(this.calibration.T1) / 8192.0) * (float64(adc) / 131072.0 - float64(this.calibration.T1) / 8192.0)) * float64(this.calibration.T3)
	t_fine := var1 + var2
	return t_fine / 5120.0, t_fine, nil
}

// Return compensated pressure in Pascals
func (this *BME280Driver) readPressure(t_fine float64) (float64, error) {
	adc, err := this.readPressureRaw()
	if err != nil {
		return 0,err
	}
	var1 := t_fine / 2.0 - 64000.0
	var2 := var1 * var1 * float64(this.calibration.P6) / 32768.0
	var2 = var2 + var1 * float64(this.calibration.P5) * 2.0
	var2 = var2 / 4.0 + float64(this.calibration.P4) * 65536.0
	var1 = (float64(this.calibration.P3) * var1 * var1 / 524288.0 + float64(this.calibration.P2) * var1) / 524288.0
	var1 = (1.0 + var1 / 32768.0) * float64(this.calibration.P1)
	if var1 == 0 {
		return 0, nil // avoid exception caused by division by zero
	}
	p := 1048576.0 - float64(adc)
	p = ((p - var2 / 4096.0) * 6250.0) / var1
	var1 = float64(this.calibration.P9) * p * p / 2147483648.0
	var2 = p * float64(this.calibration.P8) / 32768.0
	p = p + (var1 + var2 + float64(this.calibration.P7)) / 16.0
	return p / 100.0, nil
}

// Return compensated humidity
func (this *BME280Driver) readHumidity(t_fine float64) (float64, error) {
	adc, err := this.readHumidityRaw()
	if err != nil {
		return 0,err
	}
	h := t_fine - 76800.0
	h = (float64(adc) - (float64(this.calibration.H4) * 64.0 + float64(this.calibration.H5) / 16384.8 * h)) * (float64(this.calibration.H2) / 65536.0 * (1.0 + float64(this.calibration.H6) / 67108864.0 * h * (1.0 + float64(this.calibration.H3) / 67108864.0 * h)))
	h = h * (1.0 - float64(this.calibration.H1) * h / 524288.0)
	switch {
	case h > 100.0:
		return 100.0, nil
	case h < 0.0:
		return 0.0, nil
	default:
		return h, nil
	}
}

////////////////////////////////////////////////////////////////////////////////

func (this *BME280Driver) readTemperatureRaw() (int32, error) {
	msb, err := this.i2c.ReadUint8(BME280_REGISTER_TEMPDATA)
	if err != nil {
		return int32(0),err
	}
	lsb, err := this.i2c.ReadUint8(BME280_REGISTER_TEMPDATA + 1)
	if err != nil {
		return int32(0),err
	}
	xlsb, err := this.i2c.ReadUint8(BME280_REGISTER_TEMPDATA + 2)
	if err != nil {
		return int32(0),err
	}
	return ((int32(msb) << 16) | (int32(lsb) << 8) | int32(xlsb)) >> 4, nil
}

func (this *BME280Driver) readPressureRaw() (int32, error) {
	// Assumes temperature has already been read
	msb, err := this.i2c.ReadUint8(BME280_REGISTER_PRESSUREDATA)
	if err != nil {
		return int32(0),err
	}
	lsb, err := this.i2c.ReadUint8(BME280_REGISTER_PRESSUREDATA + 1)
	if err != nil {
		return int32(0),err
	}
	xlsb, err := this.i2c.ReadUint8(BME280_REGISTER_PRESSUREDATA + 2)
	if err != nil {
		return int32(0),err
	}
	return ((int32(msb) << 16) | (int32(lsb) << 8) | int32(xlsb)) >> 4, nil
}

func (this *BME280Driver) readHumidityRaw() (int32, error) {
	// Assumes temperature has already been read
	msb, err := this.i2c.ReadUint8(BME280_REGISTER_HUMIDDATA)
	if err != nil {
		return int32(0),err
	}
	lsb, err := this.i2c.ReadUint8(BME280_REGISTER_HUMIDDATA + 1)
	if err != nil {
		return int32(0),err
	}
	return (int32(msb) << 8) | int32(lsb), nil
}

