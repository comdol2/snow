package api

import (
	"fmt"
        "strings"
	"errors"
	"strconv"

)

func (c *Client) GetIncident(incidentnumber, sys_id, sysparamfields string) (interface{}, string, error) {

        if c.debug {
                fmt.Println("\nGetIncident(" + incidentnumber + ", " + sys_id + ", " + sysparamfields + ")")
        }

        if strings.HasPrefix(strings.ToUpper(incidentnumber), "INC") {

                respMap, err := c.snowTable("GET", "incident", map[string]string{"number": incidentnumber, "sysparm_fields": sysparamfields}, nil)
                if err != nil {
                        return nil, "", err
                }
                if len(respMap) < 1 {
			return respMap, "", fmt.Errorf("ERROR: No %s found", incidentnumber)
                }
                resp := respMap[0] // Actual Incident Request Record

                strResp := "Incident Number    : "          + resp.(map[string]interface{})["number"].(string)
                state := resp.(map[string]interface{})["state"].(string)
		if state == "10" {
			state = "Assigned"
		} else if state == "11" {
			state = "Work In Progress"
		} else if state == "12" {
			state = "Pending Vendor"
		} else if state == "700" {
			state = "Pending Change"
		} else if state == "13" {
			state = "Pending Customer"
		} else if state == "704" {
			state = "Pending Validation"
		} else if state == "6" {
			state = "Resolved"
		} else if state == "14" {
			state = "Canceled"
		} else {
			state = "New or Unknown"
		}
                strResp = strResp + "\nState              : " + state
                openedby := ""
                if resp.(map[string]interface{})["opened_by"] != "" {
                        openedby = resp.(map[string]interface{})["opened_by"].(map[string]interface{})["value"].(string)
                }

		opened_at := ""
		if resp.(map[string]interface{})["opened_at"] != "" {
		}
                strResp = strResp + "\nOpened At          : " + opened_at

                strResp = strResp + "\nOpened By          : " + openedby
                strResp = strResp + "\nIncident Type      : " + resp.(map[string]interface{})["u_incident_type"].(string)
                strResp = strResp + "\nContact Type       : " + resp.(map[string]interface{})["contact_type"].(string)
                cmdbci := ""
                if resp.(map[string]interface{})["cmdb_ci"] != "" {
                        cmdbci = resp.(map[string]interface{})["cmdb_ci"].(map[string]interface{})["value"].(string)
                }
                strResp = strResp + "\nConfiguration Item : " + cmdbci
                strResp = strResp + "\nCategory           : " + resp.(map[string]interface{})["category"].(string)
                strResp = strResp + "\nSubCategory        : " + resp.(map[string]interface{})["subcategory"].(string)
                assignmentgroup := ""
                if resp.(map[string]interface{})["assignment_group"] != "" {
                        assignmentgroup = resp.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string)
                }
                strResp = strResp + "\nAssignment Group   : " + assignmentgroup
                assignedto := ""
                if resp.(map[string]interface{})["assigned_to"] != "" {
                        assignedto = resp.(map[string]interface{})["assigned_to"].(map[string]interface{})["value"].(string)
                }
                strResp = strResp + "\nAssigned To        : " + assignedto
		str_impact := ""
		if resp.(map[string]interface{})["impact"] != "" {
			idx_impact, _ := strconv.Atoi(resp.(map[string]interface{})["impact"].(string))
			str_impact = strconv.Itoa(idx_impact) + " - " + arrIncidentImpactList[idx_impact-1]
		}
                strResp = strResp + "\nImpact             : " + str_impact
		str_urgency := ""
		if resp.(map[string]interface{})["urgency"] != "" {
			idx_urgency, _ := strconv.Atoi(resp.(map[string]interface{})["urgency"].(string))
			str_urgency = strconv.Itoa(idx_urgency) + " - " + arrIncidentUrgencyList[idx_urgency-1]
		}
                strResp = strResp + "\nUrgency            : " + str_urgency
		str_priority := ""
		if resp.(map[string]interface{})["priority"] != "" {
			idx_priority, _ := strconv.Atoi(resp.(map[string]interface{})["priority"].(string))
			str_priority = strconv.Itoa(idx_priority) + " - " + arrIncidentPriorityList[idx_priority-1]
		}
                strResp = strResp + "\nPriority           : " + str_priority
                strResp = strResp + "\nShort Description  : " + resp.(map[string]interface{})["short_description"].(string)
                strResp = strResp + "\nDescription        : " + strings.Replace(resp.(map[string]interface{})["description"].(string), "\n", "\n\t", -1)
                strResp = strResp + "\nWork Start Date    : " + resp.(map[string]interface{})["work_start"].(string)
                strResp = strResp + "\nWork End Date      : " + resp.(map[string]interface{})["work_end"].(string)
                closedby := ""
                if resp.(map[string]interface{})["closed_by"] != "" {
                        closedby = resp.(map[string]interface{})["closed_by"].(map[string]interface{})["value"].(string)
                }

		closed_at := ""
		if resp.(map[string]interface{})["closed_at"] != "" {
		}
                strResp = strResp + "\nClosed At          : " + closed_at

                strResp = strResp + "\nClosed By          : " + closedby
                strResp = strResp + "\nClosure Code       : " + resp.(map[string]interface{})["close_code"].(string)
                strResp = strResp + "\nClose Notes        : " + strings.Replace(resp.(map[string]interface{})["close_notes"].(string), "\n", "\n\t", -1)
                strResp = strResp + "\n"

                return resp, strResp, nil

        } else {

                return nil, "", errors.New("Provide a correct INC Number")

        }

}
