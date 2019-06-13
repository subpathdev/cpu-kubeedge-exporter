package typ

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type DeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Device `json:"items"`
}

type DeviceSpec struct {
	DeviceModelRef *v1.LocalObjectReference `json:"deviceModelRef,omitempty"`
	Protocol       ProtocolConfig           `json:"protocol,omitempty"`
	NodeSelector   *v1.NodeSelector         `json:"nodeSelector,omitempty"`
}

type ProtocolConfig struct {
	OpcUA     *ProtocolConfigOpcUA     `json:"opcua,omitempty"`
	Modbus    *ProtocolConfigModbus    `json:"modbus,omitempty"`
	Bluetooth *ProtocolConfigBluetooth `json:"bluetooth,omitempty"`
}

type ProtocolConfigOpcUA struct {
	URl            string `json:"url,omitempty"`
	UserName       string `json:"userNamem,omitempty"`
	Password       string `json:"password,omitempty"`
	SecurityPolicy string `json:"securityPolicy,omitempty"`
	SecurityMode   string `json:"securityMode,omitempty"`
	Certificate    string `json:"certificate,omitempty"`
	PrivateKey     string `json:"privateLey,omitempty"`
	Timeout        int64  `json:"timeout,omitempty"`
}

type ProtocolConfigModbus struct {
	RTU *ProtocolConfigModbusRTU `json:"rtu,omitempty"`
	TCP *ProtocolConfigModbusTCP `json:"tcp,omitempty"`
}

type ProtocolConfigModbusRTU struct {
	SerialPort string `json:"serialPort,omitempty"`
	BaudRate   int64  `json:"baudRate,omitempty"`
	DataBits   int64  `json:"dataBits,omitempty"`
	Parity     string `json:"parity,omitempty"`
	StopBits   int64  `json:"stopBits,omitempty"`
	SlaveID    int64  `json:"slaveID,omitempty"`
}

type ProtocolConfigModbusTCP struct {
	IP      string `json:"ip,omitempty"`
	Port    int64  `json:"port,omitempty"`
	SlaveID string `json:"slaveID,omitempty"`
}

type ProtocolConfigBluetooth struct {
	MACAddress string `json:"twins,omitempty"`
}

type Twin struct {
	Name    string    `json:"propertyName"`
	Actual  TwinValue `json:"reported,omitempty"`
	Desired TwinValue `json:"desired,omitempty"`
}

type TwinValue struct {
	Metadata map[string]string `json:"metadata,omitempty"`
	Value    string            `json:"value,omitempty"`
}

type DeviceStatus struct {
	Twins []Twin `json:"twins,omitempty"`
}

// Device is the Schema for the devices API
type Device struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceSpec   `json:"spec,omitempty"`
	Status DeviceStatus `json:"status,omitempty"`
}

func (in *Device) DeepCopy() *Device {
	if in == nil {
		return nil
	}

	out := new(Device)
	in.DeepCopyInto(out)
	return out
}

func (in *Device) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *Device) DeepCopyInto(out *Device) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *DeviceSpec) DeepCopyInto(out *DeviceSpec) {
	*out = *in

	if in.DeviceModelRef != nil {
		in, out := &in.DeviceModelRef, &out.DeviceModelRef
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
}

func (in *DeviceStatus) DeepCopyInto(out *DeviceStatus) {
	*out = *in
	if in.Twins != nil {
		in, out := &in.Twins, &out.Twins
		*out = make([]Twin, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *Twin) DeepCopyInto(out *Twin) {
	*out = *in
	in.Desired.DeepCopyInto(&out.Desired)
	in.Actual.DeepCopyInto(&out.Actual)
}

func (in *TwinValue) DeepCopyInto(out *TwinValue) {
	*out = *in
}

func (in *DeviceList) DeepCopyInto(out *DeviceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Device, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *DeviceList) DeepCopy() *DeviceList {
	if in == nil {
		return nil
	}
	out := new(DeviceList)
	in.DeepCopyInto(out)
	return out
}

func (in *DeviceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
