/*
Go Language Raspberry Pi Interface
(c) Copyright David Thorpe 2016
All Rights Reserved

For Licensing and Usage information, please see LICENSE.md
*/
package egl

type EGL struct {
    display Display
    context Context
    surface Surface
}

func NewContext() *EGL,error {
    this := new(EGL)

    // Initalize display
    this.display := GetDisplay()
    if err := Initialize(this.display,nil,nil); err != nil {
        return nil,err
    }
}