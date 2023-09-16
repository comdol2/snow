package api

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ===============================================
// =============== I N C I D E N T ===============
// ===============================================

var arrIncidentIncidentTypeList = []string{
	"Call Tracking",
}

var arrIncidentContactTypeList = []string{
	"Audits",
	"Automation",
	"Chat",
	"Email",
	"Governance",
	"Phone",
	"Selfservice",
	"Yammer",
}

var arrIncidentStateList = []string{
	//"New", // Should Not supported in this array
	//"Assigned", // Should Not supported in this array
	"Work In Progress",
	"Pending Vendor",
	"Pending Change",
	"Pending Customer",
	"Pending Validation",
	"Resolved",
	"Canceled",
}

var arrIncidentCloseCodesList = [][]string{
	{
		// for Work In Progress
		"",
	}, {
		// for Pending Vendor
		"",
	}, {
		// for Pending Change
		"",
	}, {
		// for Customer
		"",
	}, {
		// for Validation
		"Solved (Work Around)",
		"Solved (Permanently)",
		"Solved Remotely (Work Around)",
		"Solved Remotely (Permanently)",
		"Not Solved (Not Reproducible)",
		"Closed/Resolved by Caller",
		"Worked as Designed",
	}, {
		// for Resolved
		"Solved (Work Around)",
		"Solved (Permanently)",
		"Solved Remotely (Work Around)",
		"Solved Remotely (Permanently)",
		"Not Solved (Not Reproducible)",
		"Closed/Resolved by Caller",
		"Worked as Designed",
	}, {
		// for Cancelled
		"",
	},
}

var arrIncidentCategoryList = []string{
	"Access",
	"Error Message",
	"Performance",
	"Check",
	"Event",
	"Install",
	"Remove",
	"Repair",
	"Security",
	"Update",
}

var arrIncidentSubCategoryList = [][]string{
	{
		// for Category 'Access'
		"Data (Missing or Currupt)",
		"Hardware",
		"Network Connectivity",
		"Security Authentication / Login",
		"Software",
	}, {
		// for Category 'Error Message'
		"Hardware",
		"Job Failure",
		"Network Connectivity",
		"Security",
		"Software",
		"Storage",
	}, {
		// for Category 'Performance'
		"Degraded",
		"Function Not Working",
		"Offline",
	}, {
		// for Category 'Check'
		"Application",
		"Authentication",
		"Availability",
		"Cancellation",
		"Component",
		"Connectivity",
		"Data",
		"Functionality",
		"Performance",
		"Patch",
		"Perifheral",
		"Version",
		"Hardware",
		"Software",
		"Ticket",
		"Alert",
		"Unauthorized",
	}, {
		// for Category 'Event'
		"Application",
		"Authentication",
		"Availability",
		"Cancellation",
		"Component",
		"Connectivity",
		"Data",
		"Functionality",
		"Performance",
		"Patch",
		"Perifheral",
		"Version",
		"Hardware",
		"Software",
		"Ticket",
		"Alert",
		"Unauthorized",
	}, {
		// for Category 'Install'
		"Application",
		"Authentication",
		"Availability",
		"Cancellation",
		"Component",
		"Connectivity",
		"Data",
		"Functionality",
		"Performance",
		"Patch",
		"Perifheral",
		"Version",
		"Hardware",
		"Software",
		"Ticket",
		"Alert",
		"Unauthorized",
	}, {
		// for Category 'Remove'
		"Application",
		"Authentication",
		"Availability",
		"Cancellation",
		"Component",
		"Connectivity",
		"Data",
		"Functionality",
		"Performance",
		"Patch",
		"Perifheral",
		"Version",
		"Hardware",
		"Software",
		"Ticket",
		"Alert",
		"Unauthorized",
	}, {
		// for Category 'Repair'
		"Application",
		"Authentication",
		"Availability",
		"Cancellation",
		"Component",
		"Connectivity",
		"Data",
		"Functionality",
		"Performance",
		"Patch",
		"Perifheral",
		"Version",
		"Hardware",
		"Software",
		"Ticket",
		"Alert",
		"Unauthorized",
	}, {
		// for Category 'Security'
		"Application",
		"Authentication",
		"Availability",
		"Cancellation",
		"Component",
		"Connectivity",
		"Data",
		"Functionality",
		"Performance",
		"Patch",
		"Perifheral",
		"Version",
		"Hardware",
		"Software",
		"Ticket",
		"Alert",
		"Unauthorized",
	}, {
		// for Category 'Update'
		"Application",
		"Authentication",
		"Availability",
		"Cancellation",
		"Component",
		"Connectivity",
		"Data",
		"Functionality",
		"Performance",
		"Patch",
		"Perifheral",
		"Version",
		"Hardware",
		"Software",
		"Ticket",
		"Alert",
		"Unauthorized",
	},
}

var arrIncidentImpactList = []string{
	// Should be this order : Enterprise - Single Segment - Multiple Users - Single User
	"Enterprise",
	"Single Segment",
	"Multiple Users",
	"Single User",
}

var arrIncidentUrgencyList = []string{
	// Should be this order : Critical - High - Medium - Normal
	"Critical",
	"High",
	"Medium",
	"Normal",
}

var arrIncidentPriorityList = []string{
	// Should be this order : Critical - High - Medium - Normal
	"Critical",
	"High",
	"Medium",
	"Normal",
}

var arrIncidentCauseCodeAreaList = []string{
	"Data",
	"Hardware",
	"Software",
	"Validation",
	"Change Induced",
	"Capacity Related",
	"Knowledge/Training",
	"Security Related",
	"Component Issue",
}

var arrIncidentCauseCodeSubAreaList = [][]string{
	{
		// for Category 'Data'
		"Restored from Backup",
		"Recreted",
		"Removed",
		"Configured",
		"Data Provided",
	}, {
		// for Category 'Hardware'
		"Installed New",
		"Repaired Existing",
		"Reboot",
		"Removed",
		"Failed Over",
		"Configured",
	}, {
		// for Category 'Software'
		"Installed New",
		"Re-Installed",
		"Applied Patch",
		"Upgraded",
		"Restart",
		"Remove",
		"Licensing / Cert",
		"Configured",
	}, {
		// for Category 'Validation'
		"No Further Action Taken",
		"Withdrawn by Custoker",
		"External Vendor",
		"Customer Unreachable",
		"Test Ticket",
	}, {
		// for Category 'Change Induced'
		"IT Originated",
		"Supplier Originated",
		"User Originated",
	}, {
		// for Category 'Capacity Related'
		"Application/Licensing",
		"Local Storate",
		"Mailbox/Email",
		"Network",
		"Other",
		"PC Performance",
		"Server Performance",
		"Shared Storate",
		"Supplier",
	}, {
		// for Category 'Knowledge/Training'
		"Communication",
		"Documentation Error",
		"Research",
		"Training Needed",
	}, {
		// for Category 'Security Related'
		"Access Management",
		"Employee Breach",
		"External Attack",
		"Supplie Breach",
		"Virus/Malware",
	}, {
		// for Category 'Component Issue'
		"Environmental Issue",
		"Hardware Issue",
		"Interface Issue",
		"Software Issue",
		"Supplier Issue",
	},
}

// =================================================
// =================================================
// =================================================

func (c *Client) GetWhoAmI() string {

	user := os.Getenv("USER")

	return user

}

func (c *Client) ExistInSlice(slice []string, val string) (int, bool) {

	if c.debug {
		fmt.Println("\nExistInSlise(", slice, ", ", val, ")\n")
		fmt.Println("Slice : ", slice)
		fmt.Println("Value : ", val)
	}

	for i, item := range slice {
		if strings.ToLower(item) == strings.ToLower(strings.Replace(val, "_", " ", -1)) {
			if c.debug {
				fmt.Println("Found in above Slice!!")
			}
			return i, true
		}
	}

	if c.debug {
		fmt.Println("Not Found in above Slice!!")
	}

	// if not found
	return -1, false

}

func (c *Client) ValidateRequestedPhaseState(current, new string) error {

	if c.debug {
		fmt.Println("\nValidateRequestedPhaseState(", current, ", ", new, ")\n")
		fmt.Println("Stored Phase State : " + current)
		fmt.Println("Requested Phase State : " + new)
	}

	lowerCurrent := strings.ToLower(current)
	lowerNew := strings.ToLower(new)

	var current_index int
	var current_found bool
	if lowerCurrent == "open" {
		current_index = 2
		current_found = true
	} else {
		current_index, current_found = c.ExistInSlice(arrChangePhaseStateList, lowerCurrent)
	}
	new_index, new_found := c.ExistInSlice(arrChangePhaseStateList, lowerNew)

	if !new_found || !current_found {

		if c.debug {
			fmt.Println("Missing or wrong phase_state")
		}

		return errors.New("Missing or wrong phase_state")

	} else if current_index >= new_index {

		if c.debug {
			fmt.Println("phase_state can't be switched to previous state")
		}

		return errors.New("phase_state can't be switched to previous state")

	}

	return nil

}

func (c *Client) ValidateUserEnteredState(reqsvc, reqkind, data string, subdata map[string]string) error {

	if c.debug {
		fmt.Println("\nValidateUserEnteredState( ", reqsvc, " , ", reqkind, " , ", data, " , ", subdata, " )")
	}

	lowerReqSvc := strings.ToLower(reqsvc)
	lowerReqKind := strings.ToLower(reqkind)
	lowerData := strings.ToLower(data)

	var strOptions string
	if lowerReqSvc == "change" {

		if lowerReqKind == "category" {
			idx, cat_found := c.ExistInSlice(arrChangeCategoryList, lowerData)

			subcat_found := false
			if cat_found {
				_, subcat_found = c.ExistInSlice(arrChangeSubCategoryList[idx], strings.ToLower(subdata["subcategory"]))
			}

			if cat_found && subcat_found {
				return nil
			}

			// Getting list per error
			if !cat_found {
				strOptions = strings.Join(arrChangeCategoryList, ", ")
			} else if !subcat_found {
				lowerReqKind = "subcategory"
				strOptions = strings.Join(arrChangeSubCategoryList[idx], ", ")
			}
		} else if lowerReqKind == "state" {
			_, found := c.ExistInSlice(arrChangePhaseStateList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrChangePhaseStateList, ", ")
		} else if lowerReqKind == "closecode" {
			_, found := c.ExistInSlice(arrChangeCloseCodesList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrChangeCloseCodesList, ", ")
		}

	} else if lowerReqSvc == "ctask" {

		if lowerReqKind == "type" {
			_, found := c.ExistInSlice(arrCtaskTypeList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrCtaskTypeList, ", ")
		} else if lowerReqKind == "state" {
			idx := 0
			cat_found := false
			idx, cat_found = c.ExistInSlice(arrCtaskStateList, lowerData)

			subcat_found := false
			if cat_found {
				_, subcat_found = c.ExistInSlice(arrCtaskCloseCodesList[idx], strings.ToLower(subdata["closecode"]))
			}

			if cat_found && subcat_found {
				return nil
			}

			// Getting list per error
			if !cat_found {
				strOptions = strings.Join(arrCtaskStateList, ", ")
			} else if !subcat_found {
				lowerReqKind = "closecode"
				strOptions = strings.Join(arrCtaskCloseCodesList[idx], ", ")
			}
		}

	} else if lowerReqSvc == "incident" {

		if lowerReqKind == "incidenttype" {
			_, found := c.ExistInSlice(arrIncidentIncidentTypeList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrIncidentIncidentTypeList, ", ")
		} else if lowerReqKind == "contacttype" {
			_, found := c.ExistInSlice(arrIncidentContactTypeList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrIncidentContactTypeList, ", ")
		} else if lowerReqKind == "state" {
			state_idx, state_found := c.ExistInSlice(arrIncidentStateList, lowerData)

			closecode_found := false
			causecodearea_found := false
			causecodesubarea_found := false

			codearea_idx := -1
			_, closecode_ok := subdata["closecode"]
			_, causecodearea_ok := subdata["causecodearea"]
			_, causecodesubarea_ok := subdata["causecodesubarea"]
			if (closecode_ok && causecodearea_ok && causecodesubarea_ok) || lowerData == "pending validation" || lowerData == "resolved" {

				if state_found && len(arrIncidentCloseCodesList[state_idx]) > 0 {
					_, closecode_found = c.ExistInSlice(arrIncidentCloseCodesList[state_idx], strings.ToLower(subdata["closecode"]))
				}

				if closecode_found {
					codearea_idx, causecodearea_found = c.ExistInSlice(arrIncidentCauseCodeAreaList, strings.ToLower(subdata["causecodearea"]))
				}

				if causecodearea_found {
					_, causecodesubarea_found = c.ExistInSlice(arrIncidentCauseCodeSubAreaList[codearea_idx], strings.ToLower(subdata["causecodesubarea"]))
				}

			} else {

				closecode_found = true
				causecodearea_found = true
				causecodesubarea_found = true

			}

			if state_found && closecode_found && causecodearea_found && causecodesubarea_found {
				return nil
			}

			// Getting list per error
			if !state_found {
				strOptions = strings.Join(arrIncidentStateList, ", ")
			} else if !closecode_found {
				lowerReqKind = "closecode"
				strOptions = strings.Join(arrIncidentCloseCodesList[state_idx], ", ")
			} else if !causecodearea_found {
				lowerReqKind = "causecodearea"
				strOptions = strings.Join(arrIncidentCauseCodeAreaList, ", ")
			} else if !causecodesubarea_found {
				lowerReqKind = "causecodesubarea"
				strOptions = strings.Join(arrIncidentCauseCodeSubAreaList[codearea_idx], ", ")
			}
		} else if lowerReqKind == "category" || lowerReqKind == "subcategory" {
			idx, cat_found := c.ExistInSlice(arrIncidentCategoryList, lowerData)

			subcat_found := false
			if cat_found {
				_, subcat_found = c.ExistInSlice(arrIncidentSubCategoryList[idx], strings.ToLower(subdata["subcategory"]))
			}

			if cat_found && subcat_found {
				return nil
			}

			// Getting list per error
			if !cat_found {
				strOptions = strings.Join(arrIncidentCategoryList, ", ")
			} else if !subcat_found {
				lowerReqKind = "subcategory"
				strOptions = strings.Join(arrIncidentSubCategoryList[idx], ", ")
			}
		} else if lowerReqKind == "impact" {
			_, found := c.ExistInSlice(arrIncidentImpactList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrIncidentImpactList, ", ")
		} else if lowerReqKind == "urgency" {
			_, found := c.ExistInSlice(arrIncidentUrgencyList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrIncidentUrgencyList, ", ")
		} else if lowerReqKind == "priority" {
			_, found := c.ExistInSlice(arrIncidentPriorityList, lowerData)
			if found {
				return nil
			}
			strOptions = strings.Join(arrIncidentPriorityList, ", ")
		}

	}

	if strOptions != "" {
		strOptions = "\nYou can choose one of ( " + strOptions + " )"
	}

	return errors.New("Wrong or misspelled " + lowerReqKind + " provided" + strOptions)

}

func (c *Client) GetSysUser(sys_id, user_id string) (interface{}, string, string, error) {

	if c.debug {
		fmt.Println("\nGetSysUser( ", sys_id, ", ", user_id, " )")
	}

	if sys_id == "" && user_id == "" {
		return nil, "", "", errors.New(strErrorParamMissing)
	}

	var respSysUserMap []interface{}
	var err error
	if user_id == "" {
		respSysUserMap, err = c.SnowTable("GET", "sys_user", map[string]string{"sys_id": sys_id}, nil)
	} else {
		respSysUserMap, err = c.SnowTable("GET", "sys_user", map[string]string{"user_name": user_id}, nil)
	}
	if err != nil {
		return nil, "", "", err
	}

	var tmpUserNames, newUserNames, newSys_Id string
	for _, data := range respSysUserMap {
		if tmpUserNames == "" {
			tmpUserNames = data.(map[string]interface{})["name"].(string)
		} else {
			tmpUserNames = tmpUserNames + "|" + data.(map[string]interface{})["name"].(string)
		}
		newSys_Id = data.(map[string]interface{})["sys_id"].(string)
	}

	newUserNames = tmpUserNames

	return respSysUserMap, newUserNames, newSys_Id, nil

}

func (c *Client) GetSysUserGroup(sys_id, group_id string) (interface{}, string, error) {

	if c.debug {
		fmt.Println("\nGetSysUserGroup( ", sys_id, ", ", group_id, " )")
	}

	if sys_id == "" && group_id == "" {
		return nil, "", errors.New(strErrorParamMissing)
	}

	var respSysUserGroupMap []interface{}
	var err error
	if group_id == "" {
		respSysUserGroupMap, err = c.SnowTable("GET", "sys_user_group", map[string]string{"sys_id": sys_id}, nil)
	} else {
		respSysUserGroupMap, err = c.SnowTable("GET", "sys_user_group", map[string]string{"name": group_id}, nil)
	}
	if err != nil {
		return nil, "", err
	}

	var tmpGroupNames, newGroupNames string
	for _, data := range respSysUserGroupMap {
		if tmpGroupNames == "" {
			tmpGroupNames = data.(map[string]interface{})["name"].(string)
		} else {
			tmpGroupNames = tmpGroupNames + "|" + data.(map[string]interface{})["name"].(string)
		}
	}

	newGroupNames = tmpGroupNames

	return respSysUserGroupMap, newGroupNames, nil

}

func (c *Client) GetCmdbCi(sys_id string) (interface{}, string, string, error) {

	if c.debug {
		fmt.Println("\nGetCmdbCi( ", sys_id, " )")
	}

	if sys_id == "" {
		return nil, "", "", errors.New(strErrorParamMissing)
	}

	var respCmdbCi []interface{}
	var err error
	respCmdbCi, err = c.SnowTable("GET", "cmdb_ci", map[string]string{"sys_id": sys_id}, nil)
	if err != nil {
		return nil, "", "", err
	}

	if c.debug {
		fmt.Println(respCmdbCi)
	}

	var strCmdbCi, strEnv string
	for _, data := range respCmdbCi {
		strCmdbCi = data.(map[string]interface{})["name"].(string)
		strEnv = data.(map[string]interface{})["environment"].(string)
		if (strEnv == "") {
			strEnv = data.(map[string]interface{})["u_environment"].(string)
		}
	}

	return respCmdbCi, strCmdbCi, strEnv, nil

}

func (c *Client) GetBAPPinfo(table, name string) (interface{}, error) {

	if c.debug {
		fmt.Println("\nGetSVCinfo(", table, ", ", name, ")")
	}

	if name == "" {
		return nil, errors.New(strErrorParamMissing)
	}

	query := ""
	if table == "cmdb_ci_business_app" {
		table = "cmdb_ci_business_app"
		query = "number=" + name + "^ORnameSTARTSWITH" + name
	} else {
		table = "cmdb_ci_service_discovered"
		query = "nameSTARTSWITH" + name
	}

	respMap, err := c.SnowTable("GET", table, map[string]string{"sysparm_query": query}, nil)
	if err != nil {
		return nil, err
	}

	ci_total := len(respMap)
	ci_count := 0
	for _, resp := range respMap {
		respBappMap := resp // Actual Response
		ci_count = ci_count + 1

		fmt.Println()
		fmt.Println("[ ", ci_count, " / ", ci_total, " ]")
		fmt.Println("sys_id                      : ", respBappMap.(map[string]interface{})["sys_id"].(string))
		fmt.Println("sys_class_name              : ", respBappMap.(map[string]interface{})["sys_class_name"].(string))
		if resp.(map[string]interface{})["sys_class_name"].(string) == "cmdb_ci_business_app" {
			fmt.Println("BAPPID                      : ", respBappMap.(map[string]interface{})["number"].(string))
		} else {
			fmt.Println("Business Application sys_id : ", respBappMap.(map[string]interface{})["u_business_application"].(map[string]interface{})["value"].(string))
		}
		fmt.Println("Configuration Item          : ", respBappMap.(map[string]interface{})["name"].(string))
		fmt.Println("Environment                 : ", respBappMap.(map[string]interface{})["environment"].(string))
		fmt.Println("Life Cycle Stage            : ", respBappMap.(map[string]interface{})["life_cycle_stage"].(string))
		if respBappMap.(map[string]interface{})["cost_center"] != "" {
			fmt.Println("Cost Center                 : ", respBappMap.(map[string]interface{})["cost_center"].(map[string]interface{})["value"].(string))
		} else {
			fmt.Println("Cost Center                 :")
		}
		if respBappMap.(map[string]interface{})["owned_by"] != "" {
			_, owned_by, _, _ := c.GetSysUser(respBappMap.(map[string]interface{})["owned_by"].(map[string]interface{})["value"].(string), "")
			fmt.Println("Owned By                    : ", owned_by)
		} else {
			fmt.Println("Owned By                    :")
		}
		if respBappMap.(map[string]interface{})["u_executive_owner"] != "" {
			_, exec_owner, _, _ := c.GetSysUser(respBappMap.(map[string]interface{})["u_executive_owner"].(map[string]interface{})["value"].(string), "")
			fmt.Println("Executive Owner             : ", exec_owner)
		} else {
			fmt.Println("Executive Owner             :")
		}
		if respBappMap.(map[string]interface{})["assigned_to"] != "" {
			_, assigned_to, _, _ := c.GetSysUser(respBappMap.(map[string]interface{})["assigned_to"].(map[string]interface{})["value"].(string), "")
			fmt.Println("Assigned To                 : ", assigned_to)
		} else {
			fmt.Println("Assigned To                 :")
		}
		if respBappMap.(map[string]interface{})["support_group"] != "" {
			_, support_group, _ := c.GetSysUserGroup(respBappMap.(map[string]interface{})["support_group"].(map[string]interface{})["value"].(string), "")
			fmt.Println("Support Group               : ", support_group)
		} else {
			fmt.Println("Support Group               :")
		}
		fmt.Println("Short Description           : ", respBappMap.(map[string]interface{})["short_description"].(string))

	}

	return respMap, nil
}

func (c *Client) GetBAPPIDinfo(bappid, env string, displayOutput bool) (interface{}, string, string, string, string, string, string, string, string, string, error) {

	if c.debug {
		fmt.Println("\nGetBAPPIDinfo(", bappid, ", ", env, ", ", displayOutput, ")")
	}

	if bappid == "" {
		return nil, "", "", "", "", "", "", "", "", "", errors.New(strErrorParamMissing)
	} else {
		if env != "" {
			env = strings.Title(env)
		}
	}

	var respAppMap, respSvcMap []interface{}
	var err error
	if strings.HasPrefix(bappid, "BAPP") {

		// if bappid entered

		respAppMap, err = c.SnowTable("GET", "cmdb_ci_business_app", map[string]string{"sysparm_query": "sys_class_name=cmdb_ci_business_app^number=" + bappid}, nil)

	} else {

		// if cmdb_ci or sys_id entered

		// sys_class_name should check two tables (cmdb_ci_business_app, cmdb_ci_service_discovered) because there are legacy data which is before changemanagement refactoring.
		// we do "name=" + bappid + "^ORsys_id=" + bappid}, nil)" because we don't know exactly requested bappid is sys_id or ci_name
		respAppMap, err = c.SnowTable("GET", "cmdb_ci_business_app", map[string]string{"sysparm_query": "sys_class_name=cmdb_ci_business_app^name=" + bappid + "^ORsys_id=" + bappid}, nil)

	}

	if err != nil {
		return nil, "", "", "", "", "", "", "", "", "", err
	}
	if c.debug {
		fmt.Println(respAppMap)
	}
	if len(respAppMap) < 1 {
		return nil, "", "", "", "", "", "", "", "", "", errors.New("No Business Application found!!")
	}

	//
	// To get the bappid from the table & ci_name of Business Application
	// ==> bappid will be replaced by the actual value in cmdb_ci_business_app table
	//     because requested bappid might not be BAPPxxxxx, but sys_id or ci name
	// ==> ci_name is the name of Business Application
	//
	app_bappid := respAppMap[0].(map[string]interface{})["number"].(string)
	app_ci_name := respAppMap[0].(map[string]interface{})["name"].(string)
	app_sys_id := respAppMap[0].(map[string]interface{})["sys_id"].(string)

	// To get the exact record, env is a mandatory
	if env == "" {
		// sys_class_name should be checked in two tables (cmdb_ci_business_app, cmdb_ci_service_discovered) because there might be legacy data.
		respSvcMap, err = c.SnowTable("GET", "cmdb_ci_service_discovered", map[string]string{"sysparm_query": "sys_class_name=cmdb_ci_service_discovered^nameSTARTSWITH" + app_ci_name + "^u_business_application=" + app_sys_id}, nil)
	} else {
		// sys_class_name should be checked in two tables (cmdb_ci_business_app, cmdb_ci_service_discovered) because there might be legacy data.
		respSvcMap, err = c.SnowTable("GET", "cmdb_ci_service_discovered", map[string]string{"sysparm_query": "sys_class_name=cmdb_ci_service_discovered^environment=" + env + "^nameSTARTSWITH" + app_ci_name + "^u_business_application=" + app_sys_id}, nil)
	}

	if err != nil {
		return nil, "", "", "", "", "", "", "", "", "", err
	}
	if c.debug {
		fmt.Println(respSvcMap)
	}
	if len(respSvcMap) < 1 {
		return nil, "", "", "", "", "", "", "", "", "", errors.New("No Service Discovered!!")
	}

	var rtnBappid, rtnEnvironment, rtnCmdbCi, rtnJurisdiction, rtnCategory, rtnSubCategory, rtnAssignmentGroup, rtnPeerApprovalGroup string

	ci_total := len(respSvcMap)
	ci_count := 0

	for _, respSvc := range respSvcMap {
		resp := respSvc // Actual Response
		ci_count = ci_count + 1

		rtnBappid = app_bappid
		rtnEnvironment = resp.(map[string]interface{})["environment"].(string)
		rtnCmdbCi = resp.(map[string]interface{})["name"].(string)

		if resp.(map[string]interface{})["u_change_jurisdiction"] != "" {
			_, rtnJurisdiction, err = c.GetJurisdiction(resp.(map[string]interface{})["u_change_jurisdiction"].(map[string]interface{})["value"].(string))
			if err != nil {
				return nil, "", "", "", "", "", "", "", "", "", err
			}
		}

		if resp.(map[string]interface{})["support_group"] != "" {
			_, rtnAssignmentGroup, err = c.GetSysUserGroup(resp.(map[string]interface{})["support_group"].(map[string]interface{})["value"].(string), "")
			if err != nil {
				return nil, "", "", "", "", "", "", "", "", "", fmt.Errorf("support_group : %v", err)
			}
			if rtnAssignmentGroup == "" {
				return nil, "", "", "", "", "", "", "", "", "", fmt.Errorf("No support_group found")
			}
		}

		if resp.(map[string]interface{})["change_control"] != "" {
			_, rtnPeerApprovalGroup, err = c.GetSysUserGroup(resp.(map[string]interface{})["change_control"].(map[string]interface{})["value"].(string), "")
			if err != nil {
				return nil, "", "", "", "", "", "", "", "", "", fmt.Errorf("change_control : %v", err)
			}
			if rtnPeerApprovalGroup == "" {
				return nil, "", "", "", "", "", "", "", "", "", fmt.Errorf("No change_control found")
			}
		}

		if displayOutput || c.debug {
			fmt.Println()
			fmt.Println("[ ", ci_count, " / ", ci_total, " ]")
			if resp.(map[string]interface{})["sys_class_name"].(string) == "cmdb_ci_business_app" {
				fmt.Println("BAPPID                      : ", app_bappid)
			}
			fmt.Println("sys_id                      : ", resp.(map[string]interface{})["sys_id"].(string))
			_, exist := resp.(map[string]interface{})["u_business_application"]
			if exist {
				fmt.Println("Business Application sys_id : ", resp.(map[string]interface{})["u_business_application"].(map[string]interface{})["value"].(string))
			} else {
				fmt.Println("Business Application sys_id :")
			}
			fmt.Println("sys_class_name              : ", resp.(map[string]interface{})["sys_class_name"].(string))
			fmt.Println("Configuration Item          : ", resp.(map[string]interface{})["name"].(string))
			fmt.Println("Environment                 : ", resp.(map[string]interface{})["environment"].(string))
			fmt.Println("Life Cycle Stage            : ", resp.(map[string]interface{})["life_cycle_stage"].(string))
			fmt.Println("Change Jurisdiction         : ", rtnJurisdiction)
			fmt.Println("Assignment Group            : ", rtnAssignmentGroup)
			fmt.Println("Peer Approval Group         : ", rtnPeerApprovalGroup)
			if resp.(map[string]interface{})["cost_center"] != "" {
				fmt.Println("Cost Center                 : ", resp.(map[string]interface{})["cost_center"].(map[string]interface{})["value"].(string))
			} else {
				fmt.Println("Cost Center                 :")
			}
			if resp.(map[string]interface{})["owned_by"] != "" {
				_, owned_by, _, _ := c.GetSysUser(resp.(map[string]interface{})["owned_by"].(map[string]interface{})["value"].(string), "")
				fmt.Println("Owned By                    : ", owned_by)
			} else {
				fmt.Println("Owned By                    :")
			}
			if resp.(map[string]interface{})["u_executive_owner"] != "" {
				_, exec_owner, _, _ := c.GetSysUser(resp.(map[string]interface{})["u_executive_owner"].(map[string]interface{})["value"].(string), "")
				fmt.Println("Executive Owner             : ", exec_owner)
			} else {
				fmt.Println("Executive Owner             :")
			}
			if resp.(map[string]interface{})["assigned_to"] != "" {
				_, assigned_to, _, _ := c.GetSysUser(resp.(map[string]interface{})["assigned_to"].(map[string]interface{})["value"].(string), "")
				fmt.Println("Assigned To                 : ", assigned_to)
			} else {
				fmt.Println("Assigned To                 :")
			}
			if resp.(map[string]interface{})["support_group"] != "" {
				_, support_group, _ := c.GetSysUserGroup(resp.(map[string]interface{})["support_group"].(map[string]interface{})["value"].(string), "")
				fmt.Println("Support Group               : ", support_group)
			} else {
				fmt.Println("Support Group               :")
			}
			fmt.Println("Short Description           : ", resp.(map[string]interface{})["short_description"].(string))
		}

	}

	return respAppMap, rtnBappid, rtnEnvironment, app_sys_id, rtnCmdbCi, rtnJurisdiction, rtnCategory, rtnSubCategory, rtnAssignmentGroup, rtnPeerApprovalGroup, nil

}

func (c *Client) GetTemplateInfo(tmplsys_id, tmplname string, withdetails, displayOutput bool) (string, string, error) {

	if c.debug {
		fmt.Println("\nGetTemplateInfo(", tmplsys_id, ", ", tmplname, ", ", withdetails, ", ", displayOutput, ")")
	}

	if tmplsys_id != "" {
		tmplname = tmplsys_id
	}

	if tmplname == "" && tmplsys_id == "" {
		return "", "", errors.New("Missing the change template name or sys_id!!")
	}

	var err error
	var respTmplMap []interface{}
	var tmpl_sys_id, tmpl_name, tmpl_work_notes string
	respTmplMap, err = c.SnowTable("GET", "sys_template", map[string]string{"sysparm_query": "name=" + tmplname}, nil)
	if err != nil {
		return "", "", fmt.Errorf("Can't get the change template : %v", err)
	}

	tmpl_total := len(respTmplMap)
	tmpl_count := 0

	for _, data := range respTmplMap {
		tmpl_count = tmpl_count + 1

		tmpl_sys_id = data.(map[string]interface{})["sys_id"].(string)
		tmpl_name = data.(map[string]interface{})["name"].(string)

		tmpl_Content := strings.Replace(data.(map[string]interface{})["template"].(string), "\n", "\\n", -1)

		entries := strings.Split(data.(map[string]interface{})["template"].(string), "^")
		m := make(map[string]string)
		for _, e := range entries {
			parts := strings.Split(e, "=")
			if len(parts) > 1 {
				m[parts[0]] = parts[1]
			}
		}
		if displayOutput {
			fmt.Println("")
			fmt.Println("[ ", tmpl_count, " / ", tmpl_total, " ]")
			fmt.Println("Template Name        : ", tmpl_name)
			fmt.Println("Template ID          : ", tmpl_sys_id)
			if withdetails {
				fmt.Println("Template Type        : ", data.(map[string]interface{})["u_change_template_type"].(string))
				fmt.Println("Applied Table        : ", data.(map[string]interface{})["table"].(string))
				fmt.Println("Expiration Date      : ", data.(map[string]interface{})["u_expiration_date"].(string))
				fmt.Println("Short Description    : ", data.(map[string]interface{})["short_description"].(string))
				fmt.Println("Show On Template Bar : ", data.(map[string]interface{})["show_on_template_bar"].(string))
				fmt.Println("Template Content     : ", tmpl_Content)
				fmt.Println("")
				fmt.Println("\t---------------------------------------------------------------------------------------------")
				fmt.Println("\t SHORT DESCRIPTION : ", strings.Replace(m["short_description"], "\n", "\n\t\t", -1))
				fmt.Println("\t---------------------------------------------------------------------------------------------")
				fmt.Println("\t DESCRIPTION       : ", strings.Replace(m["description"], "\n", "\n\t\t", -1))
				fmt.Println("\t---------------------------------------------------------------------------------------------")
				fmt.Println("\t CHANGE_PLAN       : ", strings.Replace(m["change_plan"], "\n", "\n\t\t", -1))
				fmt.Println("\t---------------------------------------------------------------------------------------------")
				fmt.Println("\t TEST_PLAN         : ", strings.Replace(m["test_plan"], "\n", "\n\t\t", -1))
				fmt.Println("\t---------------------------------------------------------------------------------------------")
				fmt.Println("\t BACKOUT_PLAN      : ", strings.Replace(m["backout_plan"], "\n", "\n\t\t", -1))
				fmt.Println("\t---------------------------------------------------------------------------------------------")
				fmt.Println("\t WORK_NOTES        : ", strings.Replace(m["work_notes"], "\n", "\n\t\t", -1))
				fmt.Println("\t---------------------------------------------------------------------------------------------")
			}

		}
		if tmpl_work_notes == "" {
			tmpl_work_notes = strings.Replace(m["work_notes"], "\n", "\n\t\t", -1)
		} else {
			_, exist := data.(map[string]interface{})["work_notes"]
			if exist {
				tmpl_work_notes = tmpl_work_notes + "|" + data.(map[string]interface{})["work_notes"].(string)
			}
		}
	}

	if tmplsys_id != "" {
		return tmpl_name, tmpl_work_notes, nil
	} else {
		return tmpl_sys_id, tmpl_work_notes, nil
	}

}

func (c *Client) GetJurisdiction(sys_id string) (interface{}, string, error) {

	if c.debug {
		fmt.Println("\nGetJurisdiction( ", sys_id, " )")
	}

	if sys_id == "" {
		return nil, "", errors.New(strErrorParamMissing)
	}

	respJurisMap, err := c.SnowTable("GET", "u_change_jurisdiction", map[string]string{"sys_id": sys_id}, nil)
	if err != nil {
		return nil, "", err
	}

	var tmpChgJuris, newChgJuris string
	for _, data := range respJurisMap {
		if tmpChgJuris == "" {
			tmpChgJuris = data.(map[string]interface{})["u_name"].(string)
		} else {
			tmpChgJuris = tmpChgJuris + "|" + data.(map[string]interface{})["u_name"].(string)
		}
	}

	newChgJuris = tmpChgJuris

	return respJurisMap, newChgJuris, nil

}

func (c *Client) ConvertDateToUnixtime(inStr string) (int64, error) {

	ss := strings.Split(inStr, " ")
	if len(ss) < 2 {
		return 0, errors.New("Date format error")
	}

	dateStr := ss[0] + " " + ss[1]

	dateFormat := "2006-01-02 15:04:05"

	t, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return 0, err
	}

	outStr := t.Unix()
	return outStr, nil

}

func (c *Client) CompareBetweenTime(start_date, end_date, work_start, work_end string) error {

	if c.debug {
		fmt.Println("CompareBetweenTime(" + start_date + ", " + end_date + ", " + work_start + ", " + work_end + ")")
	}

	dateFormat := "2006-01-02 15:04:05"

	if start_date != "" && end_date != "" {

		arr_start_date := strings.Split(start_date, " ")
		arr_end_date := strings.Split(end_date, " ")
		_, err_s := time.Parse(dateFormat, arr_start_date[0]+" "+arr_start_date[1])
		_, err_e := time.Parse(dateFormat, arr_end_date[0]+" "+arr_end_date[1])
		if err_s != nil {
			return err_s
		}
		if err_e != nil {
			return err_e
		}

		start_date_unix, _ := c.ConvertDateToUnixtime(start_date)
		end_date_unix, _ := c.ConvertDateToUnixtime(end_date)

		if c.debug {
			fmt.Println("start_date_unix : ", start_date_unix)
			fmt.Println("end_date_unix   : ", end_date_unix)
		}

		if start_date_unix >= end_date_unix {
			return fmt.Errorf("Need start_date < end_date")
		}

		if work_start != "" && work_end != "" {
			work_start_unix, _ := c.ConvertDateToUnixtime(work_start)
			work_end_unix, _ := c.ConvertDateToUnixtime(work_end)

			if c.debug {
				fmt.Println("work_start_unix : ", work_start_unix)
				fmt.Println("work_end_unix   : ", work_end_unix)
			}

			if work_start_unix >= work_end_unix {
				return fmt.Errorf("Need work_start < work_end")
			}

			if work_start_unix < start_date_unix || work_end_unix < start_date_unix ||
				work_start_unix >= end_date_unix || work_start_unix > work_end_unix {
				return fmt.Errorf("work_start & work_end should be between start_date and end_date")
			}
		}

	} else if work_start != "" && work_end != "" {

		arr_work_start := strings.Split(work_start, " ")
		arr_work_end := strings.Split(work_end, " ")
		_, err_s := time.Parse(dateFormat, arr_work_start[0]+" "+arr_work_start[1])
		_, err_e := time.Parse(dateFormat, arr_work_end[0]+" "+arr_work_end[1])
		if err_s != nil {
			return err_s
		}
		if err_e != nil {
			return err_e
		}

		work_start_unix, _ := c.ConvertDateToUnixtime(work_start)
		work_end_unix, _ := c.ConvertDateToUnixtime(work_end)

		if c.debug {
			fmt.Println("work_start_unix : ", work_start_unix)
			fmt.Println("work_end_unix   : ", work_end_unix)
		}

		if work_start_unix >= work_end_unix {
			return fmt.Errorf("Need work_start < work_end")
		}

	} else {
		return fmt.Errorf("start_date & end_date & work_start & work_end not provided")
	}

	return nil

}

func (c *Client) ConvertLocalToUTC(inLocal string) (string, error) {

	inLocals := strings.Fields(inLocal)
	inLocals_len := len(inLocals)

	var parsedFormat time.Time
	var err error

	dateFormat := "2006-01-02 15:04:05"

	if inLocals_len == 2 {
		// get local timezone only
		t := time.Now()
		z, _ := t.Zone()
		reqDate := inLocal + " " + z
		parsedFormat, err = time.Parse(dateFormat+" MST", reqDate)
		if err != nil {
			return "", errors.New("Date is not in the correct format => yyyy-mm-dd hh:mm:ss ")
		}
	} else if inLocals_len == 3 {
		// check input date has offset ONLY
		parsedFormat, err = time.Parse(dateFormat+" -0700", inLocal)
		if err != nil {
			// check input date has zone ONLY
			loc := inLocals[2]
			switch loc {
			case "PST":
				loc = "America/Los_Angeles"
			case "PDT":
				loc = "America/Los_Angeles"
			case "EST":
				loc = "America/New_York"
			case "EDT":
				loc = "America/New_York"
			case "IST":
				loc = "Asia/Kolkata"
			case "JST":
				loc = "Asia/Tokyo"
			case "KST":
				loc = "Asia/Seoul"
			case "UTC":
				loc = "UTC"
			}
			reqDate := inLocals[0] + " " + inLocals[1] + " " + inLocals[2]
			tz, _ := time.LoadLocation(loc)
			parsedFormat, err := time.ParseInLocation(dateFormat+" MST", reqDate, tz)
			if err != nil {
				return "", errors.New("Date is not in the correct format with Zone")
			}
			outFormat := parsedFormat.UTC()
			UtcTime := outFormat.Format(dateFormat)

			return UtcTime, nil
		}
	} else if inLocals_len == 4 {
		// check input date has zone & offset BOTH
		parsedFormat, err = time.Parse(dateFormat+" -0700 MST", inLocal)
		if err != nil {
			return "", errors.New("Date is not in the correct format with zone & offset")
		}
	} else {
		return "", errors.New("Date is not in supported formatting")
	}

	outFormat := parsedFormat.UTC()
	UtcTime := outFormat.Format(dateFormat)

	return UtcTime, nil

}

func (c *Client) ConvertUTCToLocal(inUTC string) (string, error) {

	if c.debug {
		fmt.Println("ConvertUTCToLocal(", inUTC, ")")
	}

	layout := "2006-01-02 15:04:05 -0700 MST"
	dateFormat := "2006-01-02 15:04:05"

	t, _ := time.Parse(layout, inUTC+" +0000 UTC")

	if c.debug {
		fmt.Println("Requested UTC : ", t)
	}

	localLoc, err := time.LoadLocation("Local")
	if err != nil {
		return "", err
	}
	localDateTime := t.In(localLoc)

	convertedLocal := localDateTime.Format(dateFormat)

	if c.debug {
		fmt.Println("Converted: ", convertedLocal)
	}

	return convertedLocal, nil

}

func (c *Client) GetCurrentDateTime(hours, mins int) string {

	now := time.Now()
	if hours > 0 || mins > 0 {
		now = time.Now().Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(mins))
	}

	curDate := now.In(time.Local).Format("2006-01-02 15:04:05 MST")

	rtnDate, _ := c.ConvertLocalToUTC(curDate)

	return rtnDate

}

func (c *Client) GetUserInput(q interface{}, readLines bool) (string, error) {

	if c.debug {
		fmt.Println("\nGetUserInput(", q, ",", readLines, ")")
	}

	// if readLines = true => accept multiplines' input
	// if readLines = fause => accept ONLY a single line input

	var answer string

	scanInput := bufio.NewScanner(os.Stdin)

	if q != nil {
		var question string
		question = fmt.Sprintf("%v", q)
		fmt.Println("Question : " + question)
	}

	var lines []string
	for scanInput.Scan() {
		line := scanInput.Text()
		if readLines {
			// accept multiplines' input
			if len(line) == 1 {
				// Group Separator (GS ^]): CTRL + ]
				if line[0] == '\x1D' {
					break
				}
			}
			lines = append(lines, line)
		} else {
			// accept ONLY a single line input
			lines = append(lines, line)
			break
		}
	}

	if len(lines) > 0 {
		for _, line := range lines {
			if len(answer) == 0 {
				answer = line
			} else {
				answer = answer + "\n" + line
			}
		}
	}

	if err := scanInput.Err(); err != nil {
		return "", err
	}

	return answer, nil

}

func (c *Client) GetEnvironmentPrefix(env, upc, bappid string) string {

	NonV3ProdEnvs := []string{"live", "prod", "prod-dr", "prod-latest", "prod-load", "prod-stage", "proddr", "production"}
	V3ProdEnvs := []string{"live", "prod", "prod2", "prodtool", "prodtools", "production", "prd", "pd"}
	V3SfccProdEnvs := []string{"development", "staging", "production", "pig2-development", "pig2-staging", "pig2-production"}

	upperEnv := strings.ToUpper(env)
	upperUpc := strings.ToUpper(upc)

	lowerEnv := strings.ToLower(upperEnv)
	lowerUpc := strings.ToLower(upperUpc)

	rtnPrefix := "QA "

	if strings.HasPrefix(strings.ToLower(lowerUpc), "sfcc") {
		_, foundEnv := c.ExistInSlice(V3SfccProdEnvs, lowerEnv)
		if foundEnv {
			rtnPrefix = "PR "
		}
	} else if lowerUpc != "" {
		// V3 Pipeline
		_, foundEnv := c.ExistInSlice(V3ProdEnvs, lowerEnv)
		if foundEnv {
			rtnPrefix = "PR "
		}
	} else {
		// AWS or non V3 pipeline
		_, foundEnv := c.ExistInSlice(NonV3ProdEnvs, lowerEnv)
		if foundEnv {
			rtnPrefix = "PR "
		}
	}

	return rtnPrefix

}

func (c *Client) BuildChangeYaml(changenumber, env, bid, tname, output string, pipeline bool) (string, error) {

	if c.debug {
		fmt.Println("BuildChangeYaml(" + changenumber + ", " + env + ", " + bid + ", " + tname + ", " + output + ", " + strconv.FormatBool(pipeline) + ")")
	}

	change_request := new(Change_Request)

	// Regular expression for parsing work_notes, short_description, description, change_plan, test_plan, backout_plan, test_backout_plan
	// above 7 keys should keep below rules to be yaml.Marshal().
	// - Multiple spaces to a single space
	// - \r\n replaced to \n
	// - **** No space **** right before '\n' in parts[1]
	// - **** No space **** at the end of the string : parts[1]
	re := regexp.MustCompile(`\r?\n`)
	sp := regexp.MustCompile("  +")

	if changenumber != "" && strings.HasPrefix(strings.ToUpper(changenumber), "CHG") {

		// BuildYaml from the existing CHG ticket

		if c.debug {
			fmt.Println("***")
			fmt.Println("*** Pulling Change Request...")
			fmt.Println("***")
		}

		respChgMap, err := c.SnowTable("GET", "change_request", map[string]string{"number": changenumber, "sysparm_fields": ""}, nil)
		if err != nil {
			return "", fmt.Errorf("Can't get change_request : %v", err)
		}
		if respChgMap == nil || len(respChgMap) == 0 {
			return "", errors.New("No " + changenumber + " found")
		}
		respChg := respChgMap[0] // Actual Change Request Record

		if c.debug {
			fmt.Println(respChgMap)
		}

		change_request.Version = "1.0"

		cmdb_ci_sys_id_in_chg := ""
		if respChg.(map[string]interface{})["cmdb_ci"] != "" {
			cmdb_ci_sys_id_in_chg = respChg.(map[string]interface{})["cmdb_ci"].(map[string]interface{})["value"].(string)
		}
		_, cmdb_ci_name_in_chg, env_in_chg, _ := c.GetCmdbCi(cmdb_ci_sys_id_in_chg)

		fmt.Println("")
		fmt.Println("***")
		fmt.Println("*** Pulling The Actual CI...")
		fmt.Println("***")

		var respAppMap, respSvcMap []interface{}
		if strings.HasPrefix(strings.Title(env_in_chg), "Prod") {
			respAppMap, err = c.SnowTable("GET", "cmdb_ci_business_app", map[string]string{"sysparm_query": "sys_class_name=cmdb_ci_business_app^name=" + cmdb_ci_name_in_chg}, nil)
			if err != nil {
				return "", err
			}
		} else {
			respSvcMap, err = c.SnowTable("GET", "cmdb_ci_service_discovered", map[string]string{"sysparm_query": "sys_class_name=cmdb_ci_service_discovered^name=" + cmdb_ci_name_in_chg}, nil)
			if err != nil {
				return "", err
			}
			if respSvcMap == nil || len(respSvcMap) == 0 {
				return "", errors.New("'" + cmdb_ci_name_in_chg + "' Not found!!")
			}
			parentBappID := respSvcMap[0].(map[string]interface{})["u_business_application"].(map[string]interface{})["value"].(string)
			respAppMap, err = c.SnowTable("GET", "cmdb_ci_business_app", map[string]string{"sysparm_query": "sys_class_name=cmdb_ci_business_app^sys_id=" + parentBappID}, nil)
			if err != nil {
				return "", err
			}
		}
		cmdb_ci_name_in_chg = respAppMap[0].(map[string]interface{})["name"].(string)
		bappid_in_chg := respAppMap[0].(map[string]interface{})["number"].(string)

		if c.debug {
			fmt.Println("-------------------------------")
			fmt.Println("cmdb_ci_sys_id_in_chg : " + cmdb_ci_sys_id_in_chg)
			fmt.Println("cmdb_ci_name_in_chg : " + cmdb_ci_name_in_chg)
			fmt.Println("bappid_in_chg : " + bappid_in_chg)
			fmt.Println("env_in_chg : " + env_in_chg)
			fmt.Println("-------------------------------")
		}

		change_request.BAPPID = bappid_in_chg
		change_request.Environment = env_in_chg

		templatename := ""
		template_sys_id := ""
		if respChg.(map[string]interface{})["u_template"] != "" {
			template_sys_id = respChg.(map[string]interface{})["u_template"].(map[string]interface{})["value"].(string)
			if template_sys_id != "" {
				respTmplMap, err := c.SnowTable("GET", "sys_template", map[string]string{"sys_id": template_sys_id, "sysparm_fields": ""}, nil)
				if err != nil {
					// Should NOT block the process. Keep going no matter what.
					fmt.Println("ERROR: Can't get data from sys_template")
				}
				if respTmplMap == nil || len(respTmplMap) == 0 {
					templatename = "*** Attention: Template in " + changenumber + " might be deleted or not used anymore. Please check ***"
				} else {
					respTmpl := respTmplMap[0] // Actual Change Request Response
					templatename = respTmpl.(map[string]interface{})["name"].(string)
				}
			}
		}
		change_request.Parameters.TemplateName = templatename

		if c.debug {
			fmt.Println("templatename : " + templatename)
		}

		change_request.Parameters.AssignedTo = c.currentUser
		change_request.Parameters.RequestedBy = c.currentUser
		change_request.Parameters.Category = strings.Title(respChg.(map[string]interface{})["category"].(string))
		change_request.Parameters.SubCategory = strings.Title(respChg.(map[string]interface{})["u_subcategory"].(string))

		assignmentgroup := ""
		if respChg.(map[string]interface{})["assignment_group"] != "" {
			assignmentgroup = respChg.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string)
			if assignmentgroup != "" {
				_, assignmentgroup, _ = c.GetSysUserGroup(respChg.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string), "")
			}
		}
		change_request.Parameters.AssignmentGroup = assignmentgroup

		if c.debug {
			fmt.Println("assignmentgroup : " + assignmentgroup)
		}

		approvalgroup := ""
		if respChg.(map[string]interface{})["u_peer_approval_group"] != "" {
			approvalgroup = respChg.(map[string]interface{})["u_peer_approval_group"].(map[string]interface{})["value"].(string)
			if approvalgroup != "" {
				_, approvalgroup, _ = c.GetSysUserGroup(respChg.(map[string]interface{})["u_peer_approval_group"].(map[string]interface{})["value"].(string), "")
			}
		}
		change_request.Parameters.PeerApprovalGroup = approvalgroup

		if c.debug {
			fmt.Println("approvalgroup : " + approvalgroup)
		}

		change_request.Parameters.ShortDescription = respChg.(map[string]interface{})["short_description"].(string)
		change_request.Parameters.WorkNotes = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respChg.(map[string]interface{})["work_notes"].(string), " "), "\n"), " \n", "\n", -1), " ")
		change_request.Parameters.Description = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respChg.(map[string]interface{})["description"].(string), " "), "\n"), " \n", "\n", -1), " ")
		change_request.Parameters.ChangePlan = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respChg.(map[string]interface{})["change_plan"].(string), " "), "\n"), " \n", "\n", -1), " ")
		change_request.Parameters.TestPlan = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respChg.(map[string]interface{})["test_plan"].(string), " "), "\n"), " \n", "\n", -1), " ")
		change_request.Parameters.BackoutPlan = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respChg.(map[string]interface{})["backout_plan"].(string), " "), "\n"), " \n", "\n", -1), " ")

		if respChg.(map[string]interface{})["start_date"] != "" {
			change_request.Parameters.PlannedStartDate, _ = c.ConvertUTCToLocal(respChg.(map[string]interface{})["start_date"].(string))
		} else {
			change_request.Parameters.PlannedStartDate = ""
		}
		if respChg.(map[string]interface{})["end_date"] != "" {
			change_request.Parameters.PlannedEndDate, _ = c.ConvertUTCToLocal(respChg.(map[string]interface{})["end_date"].(string))
		} else {
			change_request.Parameters.PlannedEndDate = ""
		}

		if c.debug {
			fmt.Println("\n***")
			fmt.Println("*** Pulling Change TASKs...")
			fmt.Println("***")
		}

		respCtaskMap, err := c.SnowTable("GET", "change_task", map[string]string{"change_request": respChg.(map[string]interface{})["sys_id"].(string), "sysparm_fields": ""}, nil)
		if err != nil {
			fmt.Println("ERROR: Can't get Change Tasks of " + changenumber)
		}
		if respCtaskMap == nil {
			fmt.Println(changenumber + " doesn't have any Ctask")
		}

		if c.debug {
			fmt.Println("respCtaskMap : ", len(respCtaskMap))
		}

		// For YAML creation with CHG1234567 needs Ctasks based on CHG Ticket
		change_request.Ctasks = make([]Ctask, len(respCtaskMap))

		for i := 0; i < len(respCtaskMap); i++ {

			respCtask := respCtaskMap[i] // Actual Change Task Record

			if c.debug {
				fmt.Println(respCtask)
			}

			// ================================================================

			change_request.Ctasks[i].Type = respCtask.(map[string]interface{})["u_type"].(string)

			assignmentgroup := ""
			if respCtask.(map[string]interface{})["assignment_group"] != "" {
				assignmentgroup = respCtask.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string)
				if assignmentgroup != "" {
					_, assignmentgroup, _ = c.GetSysUserGroup(respCtask.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string), "")
				}
			}

			change_request.Ctasks[i].AssignmentGroup = assignmentgroup
			if c.debug {
				fmt.Println("assignmentgroup : " + assignmentgroup)
			}

			change_request.Ctasks[i].AssignedTo = c.currentUser
			change_request.Ctasks[i].ShortDescription = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respCtask.(map[string]interface{})["short_description"].(string), " "), "\n"), " \n", "\n", -1), " ")
			change_request.Ctasks[i].WorkNotes = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respCtask.(map[string]interface{})["work_notes"].(string), " "), "\n"), " \n", "\n", -1), " ")
			change_request.Ctasks[i].Description = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respCtask.(map[string]interface{})["description"].(string), " "), "\n"), " \n", "\n", -1), " ")
			change_request.Ctasks[i].TaskBackoutPlan = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respCtask.(map[string]interface{})["u_task_backout_plan"].(string), " "), "\n"), " \n", "\n", -1), " ")

			if respCtask.(map[string]interface{})["work_start"] != "" {
				change_request.Ctasks[i].PlannedStartDate, _ = c.ConvertUTCToLocal(respCtask.(map[string]interface{})["work_start"].(string))
			} else {
				change_request.Ctasks[i].PlannedStartDate = ""
			}

			if respCtask.(map[string]interface{})["work_end"] != "" {
				change_request.Ctasks[i].PlannedEndDate, _ = c.ConvertUTCToLocal(respCtask.(map[string]interface{})["work_end"].(string))
			} else {
				change_request.Ctasks[i].PlannedEndDate = ""
			}

			// ================================================================

		}

		_, _, respRelMap, _, err := c.GetTaskRelation("", changenumber, "", false)
		if err != nil {
			fmt.Println(err)
		} else {

			change_request.RelatedItems = make([]RelatedItem, len(respRelMap))

			for i := 0; i < len(respRelMap); i++ {
				respRel := respRelMap[i] // Actual Change Request Record
				if c.debug {
					fmt.Println("\n*** Pulling a Task....", respRel)
				}
				resp, _, err := c.GetRelatedItem(respRel.(map[string]interface{})["parent"].(map[string]interface{})["value"].(string), "", false, true, false)
				if err != nil {
					fmt.Println(err)
				}
				change_request.RelatedItems[i].Number = resp.(map[string]interface{})["number"].(string)
			}

		}

	} else {

		// Build a pure Yaml from parameters : bappid, env, template

		_, _, _, _, _, _, _, _, _, _, err := c.GetBAPPIDinfo(bid, env, false)
		if err != nil {
			if c.debug {
				fmt.Println("ERROR : ", err)
			}
			return "", errors.New("-3-Can't get Configuration Item from the provided BappID (" + bid + "). Check with --bappidcheck [bappid or ci name] or --cicheck [ci name]")
		}

		// search the exact matched template name
		rtntmplinfo, _, err3 := c.GetTemplateInfo("", tname, true, false)
		if err3 != nil {
			return "", fmt.Errorf("Can't get %s : %v", tname, err3)
		}
		if rtntmplinfo == "" {
			return "", errors.New("No '" + tname + "' template found")
		}

		respMap, err := c.SnowTable("GET", "sys_template", map[string]string{"sys_id": rtntmplinfo, "sysparm_fields": ""}, nil)
		if err != nil {
			return "", fmt.Errorf("can't get sys_template table : %v", err)
		}

		for _, data := range respMap {

			entries := strings.Split(data.(map[string]interface{})["template"].(string), "^")
			m := make(map[string]string)
			for _, e := range entries {
				parts := strings.SplitN(e, "=", 2)
				if len(parts) > 1 {
					m[parts[0]] = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(parts[1], " "), "\n"), " \n", "\n", -1), " ")
				}
			}

			change_request.Version = "1.0"
			change_request.BAPPID = bid

			change_request.RelatedItems = make([]RelatedItem, 4)
			change_request.RelatedItems[0].Number = "CHG1234567"
			change_request.RelatedItems[1].Number = "INC1234567"
			change_request.RelatedItems[2].Number = "RITM1234567"
			change_request.RelatedItems[3].Number = "PRB1234567"

			change_request.Environment = env

			change_request.Parameters.TemplateName = tname
			change_request.Parameters.AssignedTo = c.currentUser
			change_request.Parameters.RequestedBy = c.currentUser
			change_request.Parameters.AssignmentGroup = "ops-global-parks-se-guestexp"
			change_request.Parameters.PeerApprovalGroup = "ops-global-parks-se-guestexp"
			change_request.Parameters.WorkNotes = m["work_notes"]
			if env == "Prod" || env == "Production" {
				change_request.Parameters.ShortDescription = "PR "
			} else {
				change_request.Parameters.ShortDescription = "QA "
			}
			change_request.Parameters.ShortDescription = m["short_description"]
			change_request.Parameters.Description = m["description"]
			change_request.Parameters.ChangePlan = m["change_plan"]
			change_request.Parameters.TestPlan = m["test_plan"]
			change_request.Parameters.BackoutPlan = m["backout_plan"]

			change_request.Parameters.PlannedStartDate = c.GetCurrentDateTime(0, 0)
			change_request.Parameters.PlannedEndDate = c.GetCurrentDateTime(2, 0)

			// ================================================================

			if pipeline {

				change_request.BAPPID = "BAPP0199576" // This is 'WDPR PEETI OE Pipeline' cmdb

				change_request.Parameters.Category = "Business Application"
				change_request.Parameters.SubCategory = "Scheduled Release"

				// For YAML creation with BAPPID needs 4 Ctasks by default
				change_request.Ctasks = make([]Ctask, 4)

				change_request.Ctasks[0].Type = "Implementation"
				change_request.Ctasks[0].AssignmentGroup = "ops-global-peeti-oetools"
				change_request.Ctasks[0].AssignedTo = c.currentUser
				if env == "Prod" || env == "Production" {
					change_request.Ctasks[0].ShortDescription = "PR "
				} else {
					change_request.Ctasks[0].ShortDescription = "QA "
				}
				change_request.Ctasks[0].ShortDescription = change_request.Ctasks[0].ShortDescription + "WDPR GCP OE Pipeline Deploy"
				change_request.Ctasks[0].Description = "At least one or more Implementation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

				change_request.Ctasks[0].TaskBackoutPlan = "Need to define very detailed Steps to backout.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"
				change_request.Ctasks[0].PlannedStartDate = c.GetCurrentDateTime(0, 0)
				change_request.Ctasks[0].PlannedEndDate = c.GetCurrentDateTime(1, 0)

				// ================================================================

				change_request.Ctasks[1].Type = "Implementation"
				change_request.Ctasks[1].AssignmentGroup = "ops-global-peeti-oetools"
				change_request.Ctasks[1].AssignedTo = c.currentUser
				if env == "Prod" || env == "Production" {
					change_request.Ctasks[1].ShortDescription = "PR "
				} else {
					change_request.Ctasks[1].ShortDescription = "QA "
				}
				change_request.Ctasks[1].ShortDescription = change_request.Ctasks[1].ShortDescription + "WDPR GCP OE Pipeline Smoke"
				change_request.Ctasks[1].Description = "At least one or more Implementation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

				change_request.Ctasks[1].TaskBackoutPlan = "Need to define very detailed Steps to backout.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"
				change_request.Ctasks[1].PlannedStartDate = c.GetCurrentDateTime(0, 0)
				change_request.Ctasks[1].PlannedEndDate = c.GetCurrentDateTime(1, 0)

				// ================================================================

				change_request.Ctasks[2].Type = "Implementation"
				change_request.Ctasks[2].AssignmentGroup = "ops-global-peeti-oetools"
				change_request.Ctasks[2].AssignedTo = c.currentUser
				if env == "Prod" || env == "Production" {
					change_request.Ctasks[2].ShortDescription = "PR "
				} else {
					change_request.Ctasks[2].ShortDescription = "QA "
				}
				change_request.Ctasks[2].ShortDescription = change_request.Ctasks[2].ShortDescription + "WDPR GCP OE Pipeline Promote"
				change_request.Ctasks[2].Description = "At least one or more Implementation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

				change_request.Ctasks[2].TaskBackoutPlan = "Need to define very detailed Steps to backout.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"
				change_request.Ctasks[2].PlannedStartDate = c.GetCurrentDateTime(0, 0)
				change_request.Ctasks[2].PlannedEndDate = c.GetCurrentDateTime(1, 0)

				// ================================================================

				change_request.Ctasks[3].Type = "Validation"
				change_request.Ctasks[3].AssignmentGroup = m["assignment_group"]
				change_request.Ctasks[3].AssignedTo = ""
				if env == "Prod" || env == "Production" {
					change_request.Ctasks[3].ShortDescription = "PR "
				} else {
					change_request.Ctasks[3].ShortDescription = "QA "
				}
				change_request.Ctasks[3].ShortDescription = change_request.Ctasks[3].ShortDescription + change_request.Ctasks[3].Type
				change_request.Ctasks[3].Description = "At least one or more Validation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

				change_request.Ctasks[3].PlannedStartDate = c.GetCurrentDateTime(1, 0)
				change_request.Ctasks[3].PlannedEndDate = c.GetCurrentDateTime(2, 0)

			} else {

				change_request.BAPPID = bid

				change_request.Parameters.Category = "Business Application"
				change_request.Parameters.SubCategory = "General"

				// For YAML creation with BAPPID needs 2 Ctasks by default
				change_request.Ctasks = make([]Ctask, 2)

				change_request.Ctasks[0].Type = "Implementation"
				change_request.Ctasks[0].AssignmentGroup = "ops-global-peeti-oetools"
				change_request.Ctasks[0].AssignedTo = c.currentUser
				if env == "Prod" || env == "Production" {
					change_request.Ctasks[0].ShortDescription = "PR "
				} else {
					change_request.Ctasks[0].ShortDescription = "QA "
				}
				change_request.Ctasks[0].Description = "At least one or more Implementation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

				change_request.Ctasks[0].TaskBackoutPlan = "Need to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"
				change_request.Ctasks[0].PlannedStartDate = c.GetCurrentDateTime(0, 0)
				change_request.Ctasks[0].PlannedEndDate = c.GetCurrentDateTime(1, 0)

				// ================================================================

				change_request.Ctasks[1].Type = "Validation"
				change_request.Ctasks[1].AssignmentGroup = m["assignment_group"]
				change_request.Ctasks[1].AssignedTo = ""
				if env == "Prod" || env == "Production" {
					change_request.Ctasks[1].ShortDescription = "PR "
				} else {
					change_request.Ctasks[1].ShortDescription = "QA "
				}
				change_request.Ctasks[1].Description = "At least one or more Validation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

				change_request.Ctasks[1].PlannedStartDate = c.GetCurrentDateTime(1, 0)
				change_request.Ctasks[1].PlannedEndDate = c.GetCurrentDateTime(2, 0)

			}

		}

	}

	default_yaml, err4 := yaml.Marshal(change_request)
	if err4 != nil {
		log.Fatal(err4)
		return "", fmt.Errorf("Marshaling failed : %v", err4)
	}

	output, err5 := c.WriteToOutfile(output, "change", default_yaml)
	if err5 != nil {
		return "", fmt.Errorf("WriteToOutfile failed : %v", err5)
	}

	return output, nil

}

func (c *Client) WriteToOutfile(output, prefix string, content []byte) (string, error) {

	// To force the prefix to '<prefix>.xxxx' and '<prefix>.xxx.yaml'

	output_arr := strings.Split(output, "/")
	output_file := output_arr[len(output_arr)-1]
	output_file_arr := strings.Split(output_file, ".")
	output_file_ext := output_file_arr[len(output_file_arr)-1]
	if !strings.HasPrefix(strings.ToLower(output_file), prefix+".") {
		output_arr[len(output_arr)-1] = prefix + "." + output_arr[len(output_arr)-1]
	}
	if !strings.HasPrefix(strings.ToLower(output_file_ext), "yaml") {
		output_arr[len(output_arr)-1] = output_arr[len(output_arr)-1] + ".yaml"
	}
	output = strings.Join(output_arr, "/")

	// Check if output file exists
	if utils.Exists(output) {
		dt := time.Now()
		fe := fmt.Sprintf("%d%02d%02d_%02d%02d%02d", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
		backupFile := output + "." + fe
		err := os.Rename(output, backupFile)
		if err != nil {
			return output, err
		}
		fmt.Println("'" + output + "' existed but renamed!!")
	}

	err := ioutil.WriteFile(output, content, 0644)
	if err != nil {
		return output, fmt.Errorf(output+" not created : %v", err)
	}

	return output, nil

}

func (c *Client) BuildCtaskYaml(ctasknumber, output string) (string, error) {

	if c.debug {
		fmt.Println("BuildCtaskYaml(" + ctasknumber + ", " + output + ")")
	}

	ctask_request := new(Ctask_Request)

	// Regular expression for parsing work_notes, short_description, description, change_plan, test_plan, backout_plan, test_backout_plan
	// above 7 keys should keep below rules to be yaml.Marshal().
	// - Multiple spaces to a single space
	// - \r\n replaced to \n
	// - **** No space **** right before '\n' in parts[1]
	// - **** No space **** at the end of the string : parts[1]
	re := regexp.MustCompile(`\r?\n`)
	sp := regexp.MustCompile("  +")

	ctask_request.Version = "1.0"

	if ctasknumber != "" && strings.HasPrefix(strings.ToUpper(ctasknumber), "CTASK") {

		// BuildYaml from the existing CTASK ticket

		if c.debug {
			fmt.Println("***")
			fmt.Println("*** Pulling Ctask Request...")
			fmt.Println("***")
		}

		respCtaskMap, err := c.SnowTable("GET", "change_task", map[string]string{"number": ctasknumber, "sysparm_fields": ""}, nil)
		if err != nil {
			return "", errors.New("Can't get data from change_task")
		}
		if respCtaskMap == nil || len(respCtaskMap) == 0 {
			return "", errors.New("No " + ctasknumber + " found")
		}
		respCtask := respCtaskMap[0] // Actual Change Task Record

		var respChg interface{}
		change_number := ""
		if respCtask.(map[string]interface{})["change_request"] != "" {
			change_number = respCtask.(map[string]interface{})["change_request"].(map[string]interface{})["value"].(string)
			if change_number != "" {
				// change_number is sys_id at this point
				respChgMap, _ := c.SnowTable("GET", "change_request", map[string]string{"sys_id": change_number, "sysparm_fields": ""}, nil)
				if err == nil {
					// if err != nil, then skip this part
					respChg = respChgMap[0] // Actual Change Request Response
					change_number = respChg.(map[string]interface{})["number"].(string)
				}
			}
		}
		if c.debug {
			fmt.Println("change_number : " + change_number)
		}

		ctask_request.ChangeNumber = change_number

		cmdb_ci := ""
		if respChg.(map[string]interface{})["cmdb_ci"] != "" {
			cmdb_ci = respCtask.(map[string]interface{})["cmdb_ci"].(map[string]interface{})["value"].(string)
		}
		if c.debug {
			fmt.Println("cmdb_ci : " + cmdb_ci)
		}

		environment := ""
		cmdb_ci_name := ""
		if cmdb_ci != "" {
			sysparm_fields := ""
			// sys_class_name should check two tables (cmdb_ci_business_app, cmdb_ci_service_discovered) because there are legacy data which is before changemanagement refactoring.
			respCmdbCiMap, err := c.SnowTable("GET", "cmdb_ci", map[string]string{"sys_id": cmdb_ci, "sysparm_query": "sys_class_name=cmdb_ci_business_app^ORsys_class_name=cmdb_ci_service_discovered", "sysparm_fields": sysparm_fields}, nil)
			if err != nil {
				fmt.Println("ERROR: Can't get data from cmdb_ci")
			}
			if respCmdbCiMap == nil || len(respCmdbCiMap) == 0 {
				fmt.Println("*** Attention : No cmdb_ci for " + change_number + " found ***")
				cmdb_ci_name = "No cmdb_ci for " + ctasknumber + " found. Please get the correct BAPPID"
			} else {
				respCmdbCi := respCmdbCiMap[0] // Actual Change Request Response
				cmdb_ci_name = respCmdbCi.(map[string]interface{})["name"].(string)
			}
		}
		ctask_request.Environment = environment

		if c.debug {
			fmt.Println("environment : " + ctask_request.Environment)
		}

		bappid, _ := c.GetTopLevelBappid(cmdb_ci_name, ctask_request.Environment)
		ctask_request.BAPPID = bappid

		if c.debug {
			fmt.Println("bappid : " + bappid)
		}

		// For YAML creation with BAPPID needs 2 Ctasks by default
		ctask_request.Ctasks = make([]Ctask, 1)

		ctask_request.Ctasks[0].Type = respCtask.(map[string]interface{})["u_type"].(string)

		assignmentgroup := ""
		if respCtask.(map[string]interface{})["assignment_group"] != "" {
			assignmentgroup = respCtask.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string)
			if assignmentgroup != "" {
				_, assignmentgroup, _ = c.GetSysUserGroup(respCtask.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string), "")
			}
		}

		ctask_request.Ctasks[0].AssignmentGroup = assignmentgroup
		if c.debug {
			fmt.Println("assignmentgroup : " + assignmentgroup)
		}

		ctask_request.Ctasks[0].AssignedTo = c.currentUser
		ctask_request.Ctasks[0].ShortDescription = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respCtask.(map[string]interface{})["short_description"].(string), " "), "\n"), " \n", "\n", -1), " ")
		ctask_request.Ctasks[0].Description = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respCtask.(map[string]interface{})["description"].(string), " "), "\n"), " \n", "\n", -1), " ")
		ctask_request.Ctasks[0].TaskBackoutPlan = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(respCtask.(map[string]interface{})["u_task_backout_plan"].(string), " "), "\n"), " \n", "\n", -1), " ")
		ctask_request.Ctasks[0].PlannedStartDate = respCtask.(map[string]interface{})["work_start"].(string)
		ctask_request.Ctasks[0].PlannedEndDate = respCtask.(map[string]interface{})["work_end"].(string)

	} else {

		ctask_request.ChangeNumber = "CHG1234567"

		ctask_request.BAPPID = "BAPP1234567"
		ctask_request.Environment = "Latest"

		// For YAML creation with BAPPID needs 2 Ctasks by default
		ctask_request.Ctasks = make([]Ctask, 2)

		ctask_request.Ctasks[0].Type = "Implementation"
		ctask_request.Ctasks[0].AssignmentGroup = "ops-global-peeti-oetools"
		ctask_request.Ctasks[0].AssignedTo = c.currentUser
		ctask_request.Ctasks[0].ShortDescription = "PR / QA / TR / DR Implementation"
		ctask_request.Ctasks[0].Description = "At least one or more Implementation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

		ctask_request.Ctasks[0].TaskBackoutPlan = "Need to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"
		ctask_request.Ctasks[0].PlannedStartDate = c.GetCurrentDateTime(0, 0)
		ctask_request.Ctasks[0].PlannedEndDate = c.GetCurrentDateTime(1, 0)

		ctask_request.Ctasks[1].Type = "Validation"
		ctask_request.Ctasks[1].AssignmentGroup = "ops-global-peeti-oetools"
		ctask_request.Ctasks[1].AssignedTo = c.currentUser
		ctask_request.Ctasks[1].ShortDescription = "PR / QA / TR / DR Implementation"
		ctask_request.Ctasks[1].Description = "At least one or more Implementation CTASK(s) is MANDATORY!!.\nNeed to define very detailed Steps to execute.\nDescription Section of Implementation CTASK needs to collectively align with steps/scope noted in Change Plan of Parent Change under the Planning Tab"

		ctask_request.Ctasks[1].PlannedStartDate = c.GetCurrentDateTime(1, 0)
		ctask_request.Ctasks[1].PlannedEndDate = c.GetCurrentDateTime(2, 0)

	}

	default_yaml, err4 := yaml.Marshal(ctask_request)
	if err4 != nil {
		return "", fmt.Errorf("Marshaling failed : %v", err4)
	}

	output, err5 := c.WriteToOutfile(output, "ctask", default_yaml)
	if err5 != nil {
		return "", fmt.Errorf("WriteToOutfile failed : %v", err5)
	}

	return output, nil

}

func (c *Client) BuildIncidentYaml(incidentnumber, env, bid, output string) (string, error) {

	if c.debug {
		fmt.Println("BuildIncidentYaml(" + incidentnumber + ", " + env + ", " + bid + ", " + output + ")")
	}

	incident_request := new(Incident_Request)

	// Regular expression for parsing work_notes, short_description, description, change_plan, test_plan, backout_plan, test_backout_plan
	// above 7 keys should keep below rules to be yaml.Marshal().
	// - Multiple spaces to a single space
	// - \r\n replaced to \n
	// - **** No space **** right before '\n' in parts[1]
	// - **** No space **** at the end of the string : parts[1]
	re := regexp.MustCompile(`\r?\n`)
	sp := regexp.MustCompile("  +")

	if incidentnumber != "" && strings.HasPrefix(strings.ToUpper(incidentnumber), "INC") {

		// BuildYaml from the existing INC ticket

		if c.debug {
			fmt.Println("***")
			fmt.Println("*** Pulling Incident Request...")
			fmt.Println("***")
		}

		respChgMap, err := c.SnowTable("GET", "incident", map[string]string{"number": incidentnumber, "sysparm_fields": ""}, nil)
		if err != nil {
			return "", errors.New("Can't get data from incident")
		}
		if respChgMap == nil || len(respChgMap) == 0 {
			return "", errors.New("No " + incidentnumber + " found")
		}
		respChg := respChgMap[0] // Actual Change Request Record

		incident_request.Version = "1.0"

		cmdb_ci_sys_id := ""
		if respChg.(map[string]interface{})["cmdb_ci"] != "" {
			cmdb_ci_sys_id = respChg.(map[string]interface{})["cmdb_ci"].(map[string]interface{})["value"].(string)
			if c.debug {
				fmt.Println("cmdb_ci_sys_id : " + cmdb_ci_sys_id)
			}
		}

		environment := ""
		cmdb_ci_name := ""
		if cmdb_ci_sys_id != "" {
			sysparm_fields := ""
			respCmdbCiMap, err := c.SnowTable("GET", "cmdb_ci", map[string]string{"sys_id": cmdb_ci_sys_id, "sysparm_fields": sysparm_fields}, nil)
			if err != nil {
				fmt.Println("ERROR: Can't get data from cmdb_ci")
			}
			if respCmdbCiMap == nil || len(respCmdbCiMap) == 0 {
				fmt.Println("*** Attention : No cmdb_ci for " + incidentnumber + " found ***")
				cmdb_ci_name = "No cmdb_ci for " + incidentnumber + " found. Please get the correct BAPPID"
			} else {
				respCmdbCi := respCmdbCiMap[0] // Actual Change Request Response
				environment = respCmdbCi.(map[string]interface{})["environment"].(string)
				cmdb_ci_name = respCmdbCi.(map[string]interface{})["name"].(string)
			}
		}
		incident_request.Environment = environment

		if c.debug {
			fmt.Println("environment : " + incident_request.Environment)
		}

		bappid, _ := c.GetTopLevelBappid(cmdb_ci_name, incident_request.Environment)
		incident_request.BAPPID = bappid

		if c.debug {
			fmt.Println("bappid : " + bappid)
		}

		incident_request.CallerID = c.currentUser
		incident_request.OpenedBy = c.currentUser

		assignmentgroup := ""
		if respChg.(map[string]interface{})["assignment_group"] != "" {
			assignmentgroup = respChg.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string)
			if assignmentgroup != "" {
				_, assignmentgroup, _ = c.GetSysUserGroup(respChg.(map[string]interface{})["assignment_group"].(map[string]interface{})["value"].(string), "")
			}
		}
		if c.debug {
			fmt.Println("assignmentgroup : " + assignmentgroup)
		}

		incident_request.AssignmentGroup = assignmentgroup
		incident_request.AssignedTo = c.currentUser
		incident_request.Category = strings.Title(respChg.(map[string]interface{})["category"].(string))
		incident_request.SubCategory = strings.Title(respChg.(map[string]interface{})["subcategory"].(string))

		incident_request.Impact = ""
		if respChg.(map[string]interface{})["impact"] != "" {
			idx_impact, _ := strconv.Atoi(respChg.(map[string]interface{})["impact"].(string))
			incident_request.Impact = arrIncidentImpactList[idx_impact-1]
		}
		incident_request.Urgency = ""
		if respChg.(map[string]interface{})["urgency"] != "" {
			idx_urgency, _ := strconv.Atoi(respChg.(map[string]interface{})["urgency"].(string))
			incident_request.Urgency = arrIncidentUrgencyList[idx_urgency-1]
		}
		incident_request.Priority = ""
		if respChg.(map[string]interface{})["priority"] != "" {
			idx_priority, _ := strconv.Atoi(respChg.(map[string]interface{})["priority"].(string))
			incident_request.Priority = arrIncidentPriorityList[idx_priority-1]
		}

		short_description := respChg.(map[string]interface{})["short_description"].(string)
		incident_request.ShortDescription = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(short_description, " "), "\n"), " \n", "\n", -1), " ")
		description := respChg.(map[string]interface{})["description"].(string)
		incident_request.Description = strings.TrimRight(strings.Replace(re.ReplaceAllString(sp.ReplaceAllString(description, " "), "\n"), " \n", "\n", -1), " ")

	} else {
		_, _, _, _, _, _, _, _, _, _, err := c.GetBAPPIDinfo(bid, env, false)
		if err != nil {
			if c.debug {
				fmt.Println("ERROR : ", err)
			}
			return "", errors.New("-4-Can't get Configuration Item from the provided BappID (" + bid + "). Check with --bappidcheck [bappid or ci name] or --cicheck [ci name]")
		}

		incident_request.Version = "1.0"
		incident_request.BAPPID = bid
		incident_request.Environment = env

		incident_request.CallerID = c.currentUser
		incident_request.OpenedBy = c.currentUser
		incident_request.AssignmentGroup = "ops-global-parks-se-guestexp"
		incident_request.AssignedTo = c.currentUser

		incident_request.Category = ""
		incident_request.SubCategory = ""

		incident_request.Impact = ""
		incident_request.Urgency = ""
		incident_request.Priority = ""

		incident_request.ShortDescription = "this is a test short description"
		incident_request.Description = "this is a test description\n\nPlease describe in details as much as you can\n"

	}

	default_yaml, err4 := yaml.Marshal(incident_request)
	if err4 != nil {
		log.Fatal(err4)
		return "", fmt.Errorf("Marshaling failed : %v", err4)
	}

	output, err5 := c.WriteToOutfile(output, "incident", default_yaml)
	if err5 != nil {
		return "", fmt.Errorf("WriteToOutfile failed : %v", err5)
	}

	return output, nil

}

func (c *Client) VersionOrdinal(version string) string {
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}

func (c *Client) GetTopLevelBappid(cmdbci, env string) (string, error) {

	if c.debug {
		fmt.Println("*** GetTopLevelBappid(" + cmdbci + ", " + env + ")")
	}

	// if suffix has env string, then we need to remove that. i.e. " - AWS Production"
	// otherwise search bappid with the whole cmdbci
	newCmdbCi := ""
	if strings.HasSuffix(cmdbci, env) ||
		(strings.ToLower(env) == "staging" && strings.HasSuffix(cmdbci, "Stage")) ||
		(strings.ToLower(env) == "development" && strings.HasSuffix(cmdbci, "Latest")) {

		// All Configuration Item should have the format. Otherwise it will break below code
		// It should be able to be separated by " - "
		arrCmdbCi := strings.Split(cmdbci, " - ")
		lenArrCmdbCi := len(arrCmdbCi)

		if c.debug {
			for i, item := range arrCmdbCi {
				if c.debug {
					fmt.Println("arrCmdbCi[", i, "]: ", strings.Trim(item, " "))
				}
			}
		}

		// to get the cmdb_ci again without " - AWS Production"
		newLastIndexOfCmdbCi := lenArrCmdbCi - 1
		arrCmdbCi = arrCmdbCi[:newLastIndexOfCmdbCi-1]
		newCmdbCi = strings.Join(arrCmdbCi, " - ")

	} else {

		newCmdbCi = cmdbci

	}

	if c.debug {
		fmt.Println("newCmdbCi : " + newCmdbCi)
	}

	// search the exact BAPPID can be used and fields in sysparm_fields are not empty
	// TO-DO : need to add API to search bappid from sid.disney.com
	bappid := ""
	respMap, _ := c.SnowTable("GET", "cmdb_ci_business_app", map[string]string{"u_display_name": newCmdbCi}, nil)
	if c.debug {
		fmt.Println(respMap)
	}
	for _, resp := range respMap {
		// Typically BAPPID in SID.disney.com has below keys not EMPTY
		// BAPPID in your ChangeRequest YAML SHOULD start from BAPPID in SID
		if resp.(map[string]interface{})["cost_center"] == "" ||
			resp.(map[string]interface{})["owned_by"] == "" ||
			resp.(map[string]interface{})["u_executive_owner"] == "" ||
			resp.(map[string]interface{})["sys_class_name"].(string) != "cmdb_ci_business_app" {
			bappid = ""
		} else {
			bappid = resp.(map[string]interface{})["number"].(string)
		}
	}

	return bappid, nil

}

func (c *Client) SendLog(debug, dryrun bool, version, myos, snowinstance, vaultinstance, githubinstance, vaulttoken, githubtoken, egithubinstance, evaultinstance, esnowinstance, esnowusername, esnowpassword, arglist string) {
	// This is for snow cli statictics
	hostname, _ := os.Hostname()
	formData := url.Values{
		"uid":             {c.currentUser},
		"debug":           {strconv.FormatBool(debug)},
		"dryrun":          {strconv.FormatBool(dryrun)},
		"hostname":        {hostname},
		"version":         {version},
		"os":              {myos},
		"snowinstance":    {snowinstance},
		"vaultinstance":   {vaultinstance},
		"githubinstance":  {githubinstance},
		"vaulttoken":      {vaulttoken},
		"githubtoken":     {githubtoken},
		"egithubinstance": {egithubinstance},
		"evaultinstance":  {evaultinstance},
		"esnowinstance":   {esnowinstance},
		"esnowusername":   {esnowusername},
		"esnowpassword":   {esnowpassword},
		"arglist":         {arglist},
	}

	// don't care if this is successful or not
	_, err := http.PostForm("http://mon01.disid.disney.com/snow_stat.php", formData)
	if err != nil && c.debug {
		fmt.Printf("%v\n\n", err)
	}
}

var arrChgCreateList = []string{
	"category",
}
var arrChgUpdateList = []string{
	"phase",
}
var arrCtaskCreateList = []string{
	"type",
}
var arrCtaskUpdateList = []string{
	"state",
}
var arrIncCreateList = []string{
	"type",
	"contact",
	"category",
	"impact",
	"urgency",
	"priority",
}
var arrIncUpdateList = []string{
	"state",
	"causecodearea",
}

func (c *Client) ListOptions(kind, action, query string) ([]string, string, error) {

	if c.debug {

		fmt.Println("ListOptions(", kind, ", ", action, ", ", query, ")")

	}

	var arrList []string

	lowerKind := strings.ToLower(kind)
	lowerAction := strings.ToLower(action)
	lowerQuery := strings.ToLower(query)

	key := ""
	subkey := ""
	value := ""
	s := strings.Split(lowerQuery, ":")
	if len(s) == 2 {
		key = s[0]
		value = s[1]
	} else {
		key = s[0]
		value = ""
	}

	if c.debug {

		fmt.Println("lowerKiind : ", lowerKind)
		fmt.Println("lowerAction : ", lowerAction)
		fmt.Println("key : ", key)
		fmt.Println("value : ", value)

	}

	not_supported := false
	no_options_available := false

	if lowerKind == "change" {

		if lowerAction == "create" {

			if key == "category" {

				if value == "" {
					arrList = arrChangeCategoryList
				} else {

					subkey = "subcategory"

					for index, element := range arrChangeCategoryList {
						if strings.ToLower(element) == value {
							arrList = arrChangeSubCategoryList[index]
						}
					}
				}

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else {

				not_supported = true

			}

		} else if lowerAction == "update" {

			if key == "phase" {

				if value == "Completed" {
					arrList = arrChangeCloseCodesList
				} else {
					arrList = arrChangePhaseStateList
				}

				fmt.Println(len(arrList))

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else {

				not_supported = true

			}

		} else {

			not_supported = true

		}

	} else if lowerKind == "ctask" {

		if lowerAction == "create" {

			if key == "type" {

				arrList = arrCtaskTypeList

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else {

				not_supported = true

			}

		} else if lowerAction == "update" {

			if key == "state" {

				if value == "" {
					arrList = arrCtaskStateList
				} else {
					subkey = "close codes"

					for index, element := range arrCtaskStateList {
						if strings.ToLower(element) == value {
							arrList = arrCtaskCloseCodesList[index]
						}
					}
				}

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else {

				not_supported = true

			}

		} else {

			not_supported = true

		}

	} else if lowerKind == "incident" {

		if lowerAction == "create" {

			if key == "type" {

				arrList = arrIncidentIncidentTypeList

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else if key == "contact" {

				arrList = arrIncidentContactTypeList

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else if key == "impact" {

				arrList = arrIncidentImpactList

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else if key == "urgency" {

				arrList = arrIncidentUrgencyList

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else if key == "priority" {

				arrList = arrIncidentPriorityList

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else if key == "category" {

				if value == "" {
					arrList = arrIncidentCategoryList
				} else {
					subkey = "subcategory"

					for index, element := range arrIncidentCategoryList {
						if strings.ToLower(element) == value {
							arrList = arrIncidentSubCategoryList[index]
						}
					}
				}

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else {

				not_supported = true

			}

		} else if lowerAction == "update" {

			if key == "state" {
				if value == "" {
					arrList = arrIncidentStateList
				} else {
					subkey = "close codes"

					for index, element := range arrIncidentStateList {
						if strings.ToLower(element) == value {
							arrList = arrIncidentCloseCodesList[index]
						}
					}
				}

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else if key == "causecodearea" {

				if value == "" {
					arrList = arrIncidentCauseCodeAreaList
				} else {
					subkey = "cause code subarea"

					for index, element := range arrIncidentCauseCodeAreaList {
						if strings.ToLower(element) == value {
							arrList = arrIncidentCauseCodeSubAreaList[index]
						}
					}
				}

				if arrList != nil && arrList[0] == "" {
					no_options_available = true
				}

			} else {

				not_supported = true

			}

		} else {

			not_supported = true

		}

	} else {

		not_supported = true

	}

	if not_supported {
		return nil, "", fmt.Errorf("'%s' not supported", query)
	}

	if no_options_available {
		return nil, "", errors.New("No Options Available")
	}

	if subkey != "" {
		print("< Choose one of " + strings.Title(subkey) + " Options for '" + strings.Title(value) + "' below >\n\n")
	} else {
		print("< Choose one of " + strings.Title(key) + " Options below >\n\n")
	}

	strFooter := ""
	if subkey == "" {
		strFooter = "\nRun 'snow " + lowerKind + " " + lowerAction + " --listoptions \"" + key + ":<option string above>\"' to get more suboptions for each option above\n"
	}

	return arrList, strFooter, nil

}

func (c *Client) CheckState(kind, state string) bool {

	state_found := false

	if kind == "incident" {
		_, state_found = c.ExistInSlice(arrIncidentStateList, state)
	} else if kind == "ctask" {
		_, state_found = c.ExistInSlice(arrCtaskStateList, state)
	} else if kind == "change" {
		_, state_found = c.ExistInSlice(arrChangePhaseStateList, state)
	}

	return state_found

}
