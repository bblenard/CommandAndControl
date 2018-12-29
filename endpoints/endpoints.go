package endpoints

const (
	SAVETASKRESULT            string = "/Save/Task/Results"
	SAVECLIENT                string = "/Save/Clients"
	GETCOMPLETEDTASKSBYCLIENT string = "/Get/Tasks/Completed/By/Client" // Internal
	GETPENDINGTASKSBYCLIENT   string = "/Get/Tasks/Pending/By/Client"
	SAVETASK                  string = "/Save/Tasks"                 // Internal
	GETCLIENTS                string = "/Get/Clients"                // Internal
	GETCLIENTBYID             string = "/Get/Client/By/ID"           // Internal
	GETTASKBYID               string = "/Get/Task/By/ID"             // Internal
	GETTASKRESULTBYID         string = "/Get/Task/Result/By/ID"      // Internal
	GETTASKRESULTSBYCLIENT    string = "/Get/Task/Results/By/Client" // Internal
	EXPORTDB                  string = "/Export"                     // Internal
	IMPORTDB                  string = "/Import"                     // Internal
)
