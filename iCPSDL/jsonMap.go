package iCPSDL

import "fmt"

import "strings"
import "iCPSDL/util"

// var scenario1 util.JsonMap = util.JsonMap {
//     "channels": util.JsonMap {
//         "channel1": "headChannel",
//     },
//     "agents": util.JsonMap {
//         "s1": util.JsonMap {
//             "type": "headSensor",
//             "state": "tank.level",
//             "channel": "channel1",
//         },
//         "controller" : util.JsonMap  {
//             "type": "controller",
//             "inp": "channel1",
//         },
//     },
// }
//
// var scenario2 util.JsonMap = util.JsonMap {
//     "channels": util.JsonMap {
//         "channel1": "flowChannel",
//         "channel2": "flowChannel",
//         "channel3": "headChannel",
//     },
//     "agents" : util.JsonMap {
//         "s5": util.JsonMap {
//             "type": "flowSensor",
//             "state": "p1.flow",
//             "channel": "channel1",
//         },
//         "s7": util.JsonMap {
//             "type": "flowSensor",
//             "state": "p2.flow",
//             "channel": "channel2",
//         },
//         "t_est": util.JsonMap {
//             "type": "tankEstimator",
//             "height": 10.0,
//             "timestamp": 0,
//             "volume": 0.0,
//             "inflowChannel": "channel1",
//             "outflowChannel": "channel2",
//             "headChannel": "channel3",
//         },
//         "contrl": util.JsonMap {
//             "type": "controller",
//             "inp" : "channel3",
//         },
//     },
// }
//
//
// var scenario3 util.JsonMap = util.JsonMap {
//     "channels": util.JsonMap {
//         "channel1": "flowChannel",
//         "channel2": "flowChannel",
//         "channel3": "flowChannel",
//         "channel4": "headChannel",
//     },
//     "agents" : util.JsonMap {
//         "s5": util.JsonMap {
//             "type": "flowSensor",
//             "state": "p1.flow",
//             "channel": "channel1",
//         },
//         "s8": util.JsonMap {
//             "type": "flowSensor",
//             "state": "d.demand",
//             "channel": "channel2",
//         },
//         "d_est": util.JsonMap {
//             "type": "demandEstimator",
//             "demandChannel": "channel2",
//             "inflowChannel": "channel3",
//         },
//         "t_est": util.JsonMap {
//             "type": "tankEstimator",
//             "height": 10.0,
//             "timestamp": 0,
//             "volume": 0.0,
//             "inflowChannel": "channel1",
//             "outflowChannel": "channel3",
//             "headChannel": "channel4",
//         },
//         "contrl": util.JsonMap {
//             "type": "controller",
//             "inp" : "channel4",
//         },
//     },
// }
//
// var scenario4 util.JsonMap = util.JsonMap {
//     "channels": util.JsonMap {
//         "channel1": "flowChannel",
//         "channel2": "flowChannel",
//         "channel3": "flowChannel",
//         "channel4": "flowChannel",
//         "channel5": "headChannel",
//     },
//     "agents" : util.JsonMap {
//         "s2": util.JsonMap {
//             "type": "flowSensor",
//             "state": "j.demand",
//             "channel": "channel1",
//         },
//         "s4": util.JsonMap {
//             "type": "flowSensor",
//             "state": "u.flow",
//             "channel": "channel2",
//         },
//         "s7": util.JsonMap {
//             "type": "flowSensor",
//             "state": "p2.flow",
//             "channel": "channel4",
//         },
//         "j_est": util.JsonMap {
//             "type": "junctionEstimator",
//             "inflowChannel": "channel2",
//             "demandChannel": "channel1",
//             "outflowChannel": "channel3",
//         },
//         "t_est": util.JsonMap {
//             "type": "tankEstimator",
//             "height": 10.0,
//             "timestamp": 0,
//             "volume": 0.0,
//             "inflowChannel": "channel3",
//             "outflowChannel": "channel4",
//             "headChannel": "channel5",
//         },
//         "contrl": util.JsonMap {
//             "type": "controller",
//             "inp" : "channel5",
//         },
//     },
// }
//
//
// var scenario5 util.JsonMap = util.JsonMap {
//     "channels": util.JsonMap {
//         "channel1": "flowChannel",
//         "channel2": "flowChannel",
//         "channel3": "flowChannel",
//         "channel4": "flowChannel",
//         "channel5": "flowChannel",
//         "channel6": "headChannel",
//     },
//     "agents" : util.JsonMap {
//         "s2": util.JsonMap {
//             "type": "flowSensor",
//             "state": "u.flow",
//             "channel": "channel1",
//         },
//         "s4": util.JsonMap {
//             "type": "flowSensor",
//             "state": "j.demand",
//             "channel": "channel2",
//         },
//         "s8": util.JsonMap {
//             "type": "flowSensor",
//             "state": "d.demand",
//             "channel": "channel3",
//         },
//         "d_est": util.JsonMap {
//             "type": "demandEstimator",
//             "demandChannel": "channel3",
//             "inflowChannel": "channel4",
//         },
//         "j_est": util.JsonMap {
//             "type": "junctionEstimator",
//             "inflowChannel": "channel1",
//             "demandChannel": "channel2",
//             "outflowChannel": "channel5",
//         },
//         "t_est": util.JsonMap {
//             "type": "tankEstimator",
//             "height": 10.0,
//             "timestamp": 0,
//             "volume": 0.0,
//             "inflowChannel": "channel5",
//             "outflowChannel": "channel4",
//             "headChannel": "channel6",
//         },
//         "contrl": util.JsonMap {
//             "type": "controller",
//             "inp" : "channel6",
//         },
//     },
// }
//
// var scenario6 util.JsonMap = util.JsonMap {
//     "channels": util.JsonMap {
//         "channel1": "headChannel",
//         "channel2": "headChannel",
//         "channel3": "flowChannel",
//         "channel4": "flowChannel",
//         "channel5": "flowChannel",
//         "channel6": "flowChannel",
//         "channel7": "headChannel",
//     },
//     "agents" : util.JsonMap {
//         "s1": util.JsonMap {
//             "type": "headSensor",
//             "state": "r.head",
//             "channel": "channel1",
//         },
//         "s3": util.JsonMap {
//             "type": "headSensor",
//             "state": "j.head",
//             "channel": "channel2",
//         },
//         "s4": util.JsonMap {
//             "type": "flowSensor",
//             "state": "j.demand",
//             "channel": "channel3",
//         },
//         "s7": util.JsonMap {
//             "type": "flowSensor",
//             "state": "p2.flow",
//             "channel": "channel4",
//         },
//         "p_est": util.JsonMap {
//             "type": "pumpEstimator",
//             "pumpPressure": 15.0,
//             "outpressure": "channel2",
//             "inpressure": "channel1",
//             "outflow": "channel5",
//         },
//         "j_est": util.JsonMap {
//             "type": "junctionEstimator",
//             "inflowChannel": "channel5",
//             "demandChannel": "channel3",
//             "outflowChannel": "channel6",
//         },
//         "t_est": util.JsonMap {
//             "type": "tankEstimator",
//             "height": 10.0,
//             "timestamp": 0,
//             "volume": 0.0,
//             "inflowChannel": "channel6",
//             "outflowChannel": "channel4",
//             "headChannel": "channel7",
//         },
//         "contrl": util.JsonMap {
//             "type": "controller",
//             "inp" : "channel7",
//         },
//     },
// }
//
//
// var scenario7 util.JsonMap = util.JsonMap {
//     "channels": util.JsonMap {
//         "channel1": "headChannel",
//         "channel2": "headChannel",
//         "channel3": "flowChannel",
//         "channel4": "flowChannel",
//         "channel5": "flowChannel",
//         "channel6": "flowChannel",
//         "channel7": "flowChannel",
//         "channel8": "headChannel",
//     },
//     "agents" : util.JsonMap {
//         "s1": util.JsonMap {
//             "type": "headSensor",
//             "state": "r.head",
//             "channel": "channel1",
//         },
//         "s3": util.JsonMap {
//             "type": "headSensor",
//             "state": "j.head",
//             "channel": "channel2",
//         },
//         "s4": util.JsonMap {
//             "type": "flowSensor",
//             "state": "j.demand",
//             "channel": "channel3",
//         },
//         "s8": util.JsonMap {
//             "type": "flowSensor",
//             "state": "d.demand",
//             "channel": "channel4",
//         },
//         "p_est": util.JsonMap {
//             "type": "pumpEstimator",
//             "pumpPressure": 15.0,
//             "outpressure": "channel2",
//             "inpressure": "channel1",
//             "outflow": "channel5",
//         },
//         "j_est": util.JsonMap {
//             "type": "junctionEstimator",
//             "inflowChannel": "channel5",
//             "demandChannel": "channel3",
//             "outflowChannel": "channel6",
//         },
//         "d_est": util.JsonMap {
//             "type": "demandEstimator",
//             "demandChannel": "channel4",
//             "inflowChannel": "channel7",
//         },
//         "t_est": util.JsonMap {
//             "type": "tankEstimator",
//             "height": 10.0,
//             "timestamp": 0,
//             "volume": 0.0,
//             "inflowChannel": "channel6",
//             "outflowChannel": "channel7",
//             "headChannel": "channel8",
//         },
//         "contrl": util.JsonMap {
//             "type": "controller",
//             "inp" : "channel8",
//         },
//     },
// }


func (this *end) chann(participant string, channs map[string](map[string]string), counter int) int {
    return counter
}

func (this *pass) chann(participant string, channs map[string](map[string]string), counter int) int {
    return this.cont.chann(participant, channs, counter)
}

func (this *send) chann(participant string, channs map[string](map[string]string), counter int) int {
    newCounter := this.pass.chann(participant, channs, counter)
    _, ok := channs[participant]
    if ok == false {
        channs[participant] = make(map[string]string)
    }
    _, ok = channs[participant][this.to]
    if ok == false {
        _, ok = channs[this.to][participant]
        if ok == false {
            channs[participant][this.to] = fmt.Sprintf("channel%v", newCounter)
            return newCounter + 1
        } else {
            channs[participant][this.to] = channs[this.to][participant]
        }
    }

    return newCounter
}

func (this *receive) chann(participant string, channs map[string](map[string]string), counter int) int {
    newCounter := this.pass.chann(participant, channs, counter)
    _, ok := channs[participant]
    if ok == false {
        channs[participant] = make(map[string]string)
    }
    _, ok = channs[participant][this.from]
    if ok == false {
        _, ok = channs[this.from][participant]
        if ok == false {
            channs[participant][this.from] = fmt.Sprintf("channel%v", newCounter)
            return newCounter + 1
        } else {
            channs[participant][this.from] = channs[this.from][participant]
        }
    }

    return newCounter
}

func (this *loop) chann(participant string, channs map[string](map[string]string), counter int) int {
    return this.body.chann(participant, channs, counter)
}

func (this *labelVar) chann(participant string, channs map[string](map[string]string), counter int) int {
    return counter
}

func (this *choice) chann(participant string, channs map[string](map[string]string), counter int) int {
    newCounter := counter
    for _, sess := range this.sessions {
        newCounter = sess.chann(participant, channs, newCounter)
    }
    return newCounter
}

func (this *sel) chann(participant string, channs map[string](map[string]string), counter int) int {
    newCounter := this.choice.chann(participant, channs, counter)
    _, ok := channs[participant]
    if ok == false {
        channs[participant] = make(map[string]string)
    }
    _, ok = channs[participant][this.to]
    if ok == false {
        _, ok = channs[this.to][participant]
        if ok == false {
            channs[participant][this.to] = fmt.Sprintf("channel%v", newCounter)
            return newCounter + 1
        } else {
            channs[participant][this.to] = channs[this.to][participant]
        }
    }

    return newCounter
}

func (this *branch) chann(participant string, channs map[string](map[string]string), counter int) int {
    newCounter := this.choice.chann(participant, channs, counter)
    _, ok := channs[participant]
    if ok == false {
        channs[participant] = make(map[string]string)
    }
    _, ok = channs[participant][this.from]
    if ok == false {
        _, ok = channs[this.from][participant]
        if ok == false {
            channs[participant][this.from] = fmt.Sprintf("channel%v", newCounter)
            return newCounter + 1
        } else {
            channs[participant][this.from] = channs[this.from][participant]
        }
    }

    return newCounter
}


func (this *local) json() util.JsonMap {
    agents := make(util.JsonMap)
    agentChannels := make(map[string](map[string]string))
    channels := make(util.JsonMap)
    counter := 0
    for participant, sess := range this.configuration {
        agent := make(util.JsonMap)
        parts := strings.Split(participant, ".")
        counter = sess.chann(participant, agentChannels, counter)
        switch len(parts) {
            case 1:
                switch {
                    case strings.HasPrefix(participant, "controller"):
                        agent["type"] = "controller"
                        for _, chann := range agentChannels[participant] {
                            if participant != "u" {
                                agent["inp"] = chann
                            }
                            channels[chann] = "headChannel"
                        }
                    case strings.HasPrefix(participant, "u"):
                        continue
                    case strings.HasPrefix(participant, "s"):
                        if participant == "s1" {
                            agent["type"] = "headSensor"
                            agent["state"] = "r.head"
                        } else if participant == "s2" {
                            agent["type"] = "flowSensor"
                            agent["state"] = "u.flow"
                        } else if participant == "s3" {
                            agent["type"] = "headSensor"
                            agent["state"] = "j.head"
                        } else if participant == "s4" {
                            agent["type"] = "flowSensor"
                            agent["state"] = "j.demand"
                        } else if participant == "s5" {
                            agent["type"] = "flowSensor"
                            agent["state"] = "p1.flow"
                        } else if participant == "s6" {
                            agent["type"] = "headSensor"
                            agent["state"] = "tank.level"
                        } else if participant == "s7" {
                            agent["type"] = "flowSensor"
                            agent["state"] = "p2.flow"
                        } else if participant == "s8" {
                            agent["type"] = "flowSensor"
                            agent["state"] = "d.demand"
                        }
                        for _, chann := range agentChannels[participant] {
                            agent["channel"] = chann
                        }
                }
                // case for sensor, controller or actuator
            case 2:
                // case for estimators
                switch parts[1] {
                    case "tank_mass":
                        agent["type"] = "tankEstimator"
                        for adv, chann := range agentChannels[participant] {
                            switch {
                                case adv =="controller":
                                    agent["headChannel"] = chann
                                case adv == "s5" || adv == "j.junction_mass":
                                    agent["inflowChannel"] = chann
                                    channels[chann] = "flowChannel"
                                case adv == "s7" || adv == "d.demand_mass":
                                    agent["outflowChannel"] = chann
                                    channels[chann] = "flowChannel"
                            }
                        }
                    case "junction_mass":
                        agent["type"] = "junctionEstimator"
                        for adv, chann := range agentChannels[participant] {
                            switch {
                                case adv == "s4" || adv == "u.link_energy":
                                    agent["inflowChannel"] = chann
                                    channels[chann] = "flowChannel"
                                case adv == "s2":
                                    agent["demandChannel"] = chann
                                    channels[chann] = "flowChannel"
                                default:
                                    agent["outflowChannel"] = chann
                                    channels[chann] = "flowChannel"
                            }
                        }
                    case "demand_mass":
                        agent["type"] = "demandEstimator"
                        for adv, chann := range agentChannels[participant] {
                            switch {
                                case adv =="s8":
                                    agent["demandChannel"] = chann
                                    channels[chann] = "flowChannel"
                                default:
                                    agent["inflowChannel"] = chann
                                    channels[chann] = "flowChannel"
                            }
                        }
                    case "link_energy":
                        agent["type"] = "pumpEstimator"
                        for adv, chann := range agentChannels[participant] {
                            switch {
                                case adv =="s1":
                                    agent["inpressure"] = chann
                                    channels[chann] = "headChannel"
                                case adv == "s3":
                                    agent["outpressure"] = chann
                                    channels[chann] = "headChannel"
                                default:
                                    agent["outflow"] = chann
                                    channels[chann] = "flowChannel"
                            }
                        }
                }
        }
        agents[participant] = agent
    }
    scenario := make(util.JsonMap)
    scenario["channels"] = channels
    scenario["agents"] = agents
    //fmt.Println(channels)
    return scenario
}

// func (this *local) json() util.JsonMap {
//     switch len(this.configuration) {
//         case 3:
//             return scenario1
//         case 5:
//             return scenario2
//         case 6:
//             return scenario3
//         case 7:
//             return scenario4
//         case 8:
//             return scenario5
//         case 9:
//             return scenario6
//         case 10:
//             return scenario7
//     }
//     return scenario1
// }
