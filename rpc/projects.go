package rpc

import (
	"fmt"
	"net/http"

	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

// collect all handlers in one place so that we can aseemble them easily
// there are some repeating RPCs that we would like to avoid and maybe there's some
// nice way to group them together
// setupProjectRPCs sets up all the RPC calls related to projects that might be used
func setupProjectRPCs() {
	insertProject()
	getProject()
	getAllProjects()
	getPreOriginProjects()
	getOriginProjects()
	getProposedProjects()
	getFinalProjects()
	getFundedProjects()
}

func parseProject(r *http.Request) (solar.Project, error) {
	// we need to create an instance of the Project
	// and then map values if they do exist
	// note that we just prepare the project here and don't invest in it
	// for that, we need new a new investor struct and a recipient struct
	var prepProject solar.Project
	err := r.ParseForm()
	if err != nil {
		return prepProject, err
	}
	// if we're inserting this in, we need to get the next index number
	// so that we can set this without causing some weird bugs
	allProjects, err := solar.RetrieveAllProjects()
	if err != nil {
		return prepProject, fmt.Errorf("Error in assigning index")
	}
	prepProject.Params.Index = len(allProjects) + 1
	if r.FormValue("PanelSize") == "" || r.FormValue("TotalValue") == "" || r.FormValue("Location") == "" || r.FormValue("Metadata") == "" || r.FormValue("Stage") == "" {
		return prepProject, fmt.Errorf("One of given params is missing: PanelSize, TotalValue, Location, Metadata")
	}

	prepProject.Params.PanelSize = r.FormValue("PanelSize")
	prepProject.Params.TotalValue = utils.StoI(r.FormValue("TotalValue"))
	prepProject.Params.Location = r.FormValue("Location")
	prepProject.Params.Metadata = r.FormValue("Metadata")
	prepProject.Stage = utils.StoF(r.FormValue("Stage"))
	prepProject.Params.MoneyRaised = 0
	prepProject.Params.BalLeft = float64(0)
	prepProject.Params.DateInitiated = utils.Timestamp()
	return prepProject, nil
}

func insertProject() {
	// this should be a post method since you want to accept an project and then insert
	// that into the database
	// this route does not define an originator and would mostly not be useful, should
	// look into a way where we can define originators in the route as well
	http.HandleFunc("/project/insert", func(w http.ResponseWriter, r *http.Request) {
		checkPost(w, r)
		var prepProject solar.Project
		prepProject, err := parseProject(r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		err = prepProject.Save()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		var rt StatusResponse
		rt.Status = 200
		MarshalSend(w, r, rt)
	})
}

func getAllProjects() {
	http.HandleFunc("/project/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// make a call to the db to get all projects
		// while making this call, the rpc should not be aware of the db we are using
		// and stuff. So we need to have another route that would open the existing
		// db, without asking for one
		allProjects, err := solar.RetrieveAllProjects()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		MarshalSend(w, r, allProjects)
	})
}

func getProject() {
	// we need to read passed the key from the URL that the user calls
	http.HandleFunc("/project/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		if r.URL.Query()["index"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		contract, err := solar.RetrieveProject(uKey)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		MarshalSend(w, r, contract)
	})
}

func projectHandler(w http.ResponseWriter, r *http.Request, stage float64) {
	checkGet(w, r)
	allProjects, err := solar.RetrieveProjectsAtStage(stage)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	MarshalSend(w, r, allProjects)
}

func getPreOriginProjects() {
	http.HandleFunc("/project/preorigin", func(w http.ResponseWriter, r *http.Request) {
		projectHandler(w, r, 0)
	})
}

func getOriginProjects() {
	http.HandleFunc("/project/origin", func(w http.ResponseWriter, r *http.Request) {
		projectHandler(w, r, 1)
	})
}

func getProposedProjects() {
	// we need to read passed the key from the URL that the user calls
	http.HandleFunc("/project/proposed", func(w http.ResponseWriter, r *http.Request) {
		projectHandler(w, r, 2)
	})
}

func getFinalProjects() {
	// we need to read passed the key from the URL that the user calls
	http.HandleFunc("/project/final", func(w http.ResponseWriter, r *http.Request) {
		projectHandler(w, r, 3)
	})
}

func getFundedProjects() {
	// we need to read passed the key from the URL that the user calls
	http.HandleFunc("/project/funded", func(w http.ResponseWriter, r *http.Request) {
		projectHandler(w, r, 4)
	})
}
