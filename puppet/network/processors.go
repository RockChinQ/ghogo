package network

import (
	"ghogo/puppet/operation/subprocess"
	"ghogo/util/puppet"
	"time"
)

func (p *PuppetPayloadSubProcess) Process(h *ConsoleHandler) {
	switch p.Operation {
	case "RUN":
		_, err := subprocess.MakeProcessHandler(p.SubProcess, p.Content, p.Args, p.Decoding, func(s string, sb subprocess.ProcessHandler) {

			outpayload := PuppetPayloadSubProcessOut{
				puppet.PayloadSubProcessOut{
					SubProcess: sb.UUID,
					Status:     "ALIVE",
					Content:    sb.Buffer,
					TimeStamp:  time.Now().Unix(),
				},
			}
			HandlerInst.IO.AppendNetPackage("PAYLOAD_SUB_PROCESS_OUT", outpayload)
		})
		if err != nil {
			p.Result = "fail:" + err.Error()
		} else {
			p.Result = "succ"
		}

		HandlerInst.IO.AppendNetPackage("PAYLOAD_SUB_PROCESS", *p)

	}
}
