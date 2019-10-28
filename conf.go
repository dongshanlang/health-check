package conf

import "strconv"

type Conf struct {
	IP     string
	Port   int
	custom bool
}

func (c *Conf) Addr() string {
	if c.custom {
		return c.IP + strconv.Itoa(c.Port)
	}
	return ":8080"
}
func (c *Conf) SetIP(ip string) error {
	c.IP = ip
	return nil
}
func (c *Conf) SetPort(port int) error {
	c.Port = port
	return nil
}
