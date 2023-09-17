package api

import (
	"os"
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

