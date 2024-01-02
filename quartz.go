package quartz

import (
	"fmt"
	"io"
	"net"
)

type Quartz struct {
	address    string
	port       uint16
	magnum     bool
	conn       net.Conn
	rxBuff     []byte
	rxBuffLen  uint
	RxMessages chan QuartzResponse
	end        chan bool
}

func NewQuartz(address string, port uint16, isMagnum bool) *Quartz {
	return &Quartz{
		address:    address,
		port:       port,
		magnum:     isMagnum,
		conn:       nil,
		rxBuff:     make([]byte, 256),
		rxBuffLen:  0,
		RxMessages: make(chan QuartzResponse, 100),
		end:        make(chan bool),
	}
}

func (q *Quartz) Connect() error {
	address := fmt.Sprintf("%s:%d", q.address, q.port)
	c, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	q.conn = c
	go q.rxLoop()
	return nil
}

func (q *Quartz) Disconnect() error {
	err := q.conn.Close()
	// Ignore if we are already disconnected
	if err == io.EOF {
		err = nil
	}
	q.conn = nil
	return err
}

func (q *Quartz) rxLoop() {
	for {
		select {
		case <-q.end:
			return
		default:
			q.updateBuffer()
			q.updateRxLines()
		}
	}
}

func (q *Quartz) updateBuffer() {
	if q.conn == nil {
		return
	}

	buff := make([]byte, 256)
	n, err := q.conn.Read(buff)
	if err == io.EOF {
		q.Disconnect()
		return
	}
	for i := 0; i < n; i++ {
		q.rxBuff[q.rxBuffLen] = buff[i]
		q.rxBuffLen++
	}
}

func (q *Quartz) updateRxLines() {
	start := -1
	lastEnd := -1
	for i := 0; i < int(q.rxBuffLen); i++ {
		if start == -1 && string(q.rxBuff[i]) == "." {
			start = i
		} else if start != -1 {
			if string(q.rxBuff[i]) == "\r" {
				// Full string found
				str := string(q.rxBuff[start : i+1])
				q.RxMessages <- parseResponse(str)
				start = -1
				lastEnd = i
			}
		}
	}
	// Clear buffer
	newBuff := q.rxBuff[lastEnd+1:]
	n := copy(q.rxBuff, newBuff)
	q.rxBuffLen = uint(n)
}

func (q *Quartz) sendCommand(cmd string) error {
	cmdBytes := []byte(cmd)
	_, err := q.conn.Write(cmdBytes)
	if err != nil {
		if err == io.EOF {
			q.Disconnect()
		} else {
			return err
		}
	}

	return nil
}

func (q *Quartz) SetCrosspoint(levels []QuartzLevel, dest uint, src uint) error {
	level := sortLevelsToString(levels)
	// Send message (.S{LEVEL}{DEST},{SOURCE}(CR))
	msg := fmt.Sprintf(".S%s%d,%d\r", level, dest, src)
	return q.sendCommand(msg)
}

func (q *Quartz) LockDestination(dest uint) error {
	msg := fmt.Sprintf(".BL%d\r", dest)
	return q.sendCommand(msg)
}

func (q *Quartz) UnlockDestination(dest uint) error {
	msg := fmt.Sprintf(".BU%d\r", dest)
	return q.sendCommand(msg)
}

func (q *Quartz) GetDestinationLock(dest uint) error {
	msg := fmt.Sprintf(".BI%d\r", dest)
	return q.sendCommand(msg)
}

func (q *Quartz) FireSystemSalvo(salvoId int) error {
	if salvoId > 999 {
		return fmt.Errorf("salvoId must be at most 3 digits")
	}
	msg := fmt.Sprintf(".F%d\r", salvoId)
	return q.sendCommand(msg)
}

func (q *Quartz) GetRoute(level QuartzLevel, dest uint) error {
	msg := fmt.Sprintf(".I%s%d\r", level, dest)
	return q.sendCommand(msg)
}

func (q *Quartz) GetDestinationName(dest uint) error {
	msg := fmt.Sprintf(".RD%d\r", dest)
	return q.sendCommand(msg)
}

func (q *Quartz) GetSourceName(src uint) error {
	msg := fmt.Sprintf(".RS%d\r", src)
	return q.sendCommand(msg)
}

func (q *Quartz) GetLevelName(level QuartzLevel) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".RL%s\r", level)
	return q.sendCommand(msg)
}

func (q *Quartz) GetDestinationButtonName(dest uint) error {
	msg := fmt.Sprintf(".RE%d\r", dest)
	return q.sendCommand(msg)
}

func (q *Quartz) GetSourceButtonName(src uint) error {
	msg := fmt.Sprintf(".RT%d\r", src)
	return q.sendCommand(msg)
}

func (q *Quartz) GetLevelButtonName(level QuartzLevel) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".RM%s\r", level)
	return q.sendCommand(msg)
}

func (q *Quartz) WriteDestinationName(dest uint, name string) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	if len(name) > 8 {
		return fmt.Errorf("name (%s) must be less than 8 characters (is %d)", name, len(name))
	}
	msg := fmt.Sprintf(".WD%d,%s\r", dest, name)
	return q.sendCommand(msg)
}

func (q *Quartz) WriteSourceName(src uint, name string) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	if len(name) > 8 {
		return fmt.Errorf("name (%s) must be less than 8 characters (is %d)", name, len(name))
	}
	msg := fmt.Sprintf(".WS%d,%s\r", src, name)
	return q.sendCommand(msg)
}

func (q *Quartz) WriteLevelName(level QuartzLevel, name string) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	if len(name) > 8 {
		return fmt.Errorf("name (%s) must be less than 8 characters (is %d)", name, len(name))
	}
	msg := fmt.Sprintf(".WS%s,%s\r", level, name)
	return q.sendCommand(msg)
}

func (q *Quartz) WriteDestinationButtonName(dest uint, name string) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	if len(name) > 10 {
		return fmt.Errorf("name (%s) must be less than 10 characters (is %d)", name, len(name))
	}
	msg := fmt.Sprintf(".WE%d,%s\r", dest, name)
	return q.sendCommand(msg)
}

func (q *Quartz) WriteSourceButtonName(src uint, name string) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	if len(name) > 10 {
		return fmt.Errorf("name (%s) must be less than 10 characters (is %d)", name, len(name))
	}
	msg := fmt.Sprintf(".WT%d,%s\r", src, name)
	return q.sendCommand(msg)
}

func (q *Quartz) WriteLevelButtonName(level QuartzLevel, name string) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	if len(name) > 10 {
		return fmt.Errorf("name (%s) must be less than 10 characters (is %d)", name, len(name))
	}
	msg := fmt.Sprintf(".WM%s,%s\r", level, name)
	return q.sendCommand(msg)
}

func (q *Quartz) SalvoSelect(salvoId uint) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".QC%d\r", salvoId)
	return q.sendCommand(msg)
}

func (q *Quartz) SalvoEmpty(salvoId uint) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".QR%d\r", salvoId)
	return q.sendCommand(msg)
}

func (q *Quartz) SalvoAddCrosspoint(levels []QuartzLevel, dest uint, src uint) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	level := sortLevelsToString(levels)
	// Send message (.S{LEVEL}{DEST},{SOURCE}(CR))
	msg := fmt.Sprintf(".QS%s%d,%d\r", level, dest, src)
	return q.sendCommand(msg)
}

func (q *Quartz) SalvoFireNow(salvoId uint) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".QF%d\r", salvoId)
	return q.sendCommand(msg)
}

func (q *Quartz) SalvoFireAtTime(salvoId uint, hour uint, minute uint, second uint, frame uint) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".QF%dT1:%02d:%02d:%02d:%02d\r", salvoId, hour, minute, second, frame)
	return q.sendCommand(msg)
}

func (q *Quartz) SalvoDelete(salvoId uint) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".QD%d\r", salvoId)
	return q.sendCommand(msg)
}

func (q *Quartz) SalvoListItemCount(salvoId uint) error {
	if q.magnum {
		return fmt.Errorf("command not supported by magnum")
	}
	msg := fmt.Sprintf(".QL%d\r", salvoId)
	return q.sendCommand(msg)
}

func (q *Quartz) SendPing() error {
	msg := ".#01\r"
	return q.sendCommand(msg)
}
