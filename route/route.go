package route

import (
	"net/http"

	"login-service/config"
	"login-service/controller"
	"login-service/helper/at"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}
	config.SetEnv()

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		controller.GetHome(w, r)
	//jalan setiap jam 3 pagi
	case method == "GET" && path == "/refresh/token":
		controller.GetNewToken(w, r)
	case method == "GET" && path == "/data/pushrepo/kemarin":
		controller.GetYesterdayDistincWAGroup(w, r)
	case method == "GET" && path == "/data/user":
		controller.GetDataUser(w, r)
	//generate token akses untuk kirim wa jangka panjang token
	case method == "PUT" && path == "/data/user":
		controller.PutTokenDataUser(w, r)
	case method == "GET" && path == "/data/user/task/todo":
		controller.GetTaskUser(w, r)
	case method == "GET" && path == "/data/user/task/doing":
		controller.GetTaskDoing(w, r)
	case method == "PUT" && path == "/data/user/task/doing":
		controller.PutTaskUser(w, r)
	case method == "GET" && path == "/data/user/task/done":
		controller.GetTaskDone(w, r)
	case method == "POST" && path == "/data/user/task/done":
		controller.PostTaskUser(w, r)
	case method == "POST" && path == "/data/user":
		controller.PostDataUser(w, r)
	case method == "POST" && at.URLParam(path, "/data/user/wa/:nomorwa"):
		controller.PostDataUserFromWA(w, r)
	case method == "POST" && path == "/data/proyek":
		controller.PostDataProject(w, r)
	case method == "GET" && path == "/data/proyek":
		controller.GetDataProject(w, r)
	case method == "PUT" && path == "/data/proyek":
		controller.PutDataProject(w, r)
	case method == "DELETE" && path == "/data/proyek":
		controller.DeleteDataProject(w, r)
	case method == "GET" && path == "/data/proyek/anggota":
		controller.GetDataMemberProject(w, r)
	case method == "POST" && path == "/data/proyek/anggota":
		controller.PostDataMemberProject(w, r)
	case method == "POST" && path == "/approvebimbingan":
		controller.ApproveBimbinganbyPoin(w, r)
	case method == "DELETE" && path == "/data/proyek/anggota":
		controller.DeleteDataMemberProject(w, r)
	case method == "POST" && at.URLParam(path, "/webhook/github/:proyek"):
		controller.PostWebHookGithub(w, r)
	case method == "POST" && at.URLParam(path, "/webhook/gitlab/:proyek"):
		controller.PostWebHookGitlab(w, r)
	case method == "POST" && path == "/notif/ux/postlaporan":
		controller.PostLaporan(w, r)
	case method == "POST" && path == "/notif/ux/postfeedback": //posting feedback
		controller.PostFeedback(w, r)
	case method == "POST" && path == "/notif/ux/postrating": //resume atau risalah rapat dan feedback
		controller.PostRatingLaporan(w, r)
	case method == "POST" && path == "/notif/ux/postmeeting":
		controller.PostMeeting(w, r)
	case method == "POST" && at.URLParam(path, "/notif/ux/postpresensi/:id"):
		controller.PostPresensi(w, r)
	case method == "POST" && at.URLParam(path, "/notif/ux/posttasklists/:id"):
		controller.PostTaskList(w, r)
	case method == "GET" && at.URLParam(path, "/notif/ux/getlaporan/:id"):
		controller.GetLaporan(w, r)
	case method == "POST" && at.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)
	// LMS
	case method == "GET" && path == "/lms/refresh/cookie":
		controller.RefreshLMSCookie(w, r)
	case method == "GET" && path == "/lms/count/user":
		controller.GetCountDocUser(w, r)
	// Google Auth
	case method == "POST" && path == "/auth/users":
		controller.Auth(w, r)
	case method == "POST" && path == "/auth/login":
		controller.GeneratePasswordHandler(w, r)
	case method == "POST" && path == "/auth/verify":
		controller.VerifyPasswordHandler(w, r)
	case method == "POST" && path == "/auth/resend":
		controller.ResendPasswordHandler(w, r)
	// Google Auth
	// karyawan
	case method == "GET" && path == "/data/karyawan":
		controller.GetDataBio(w, r)
	default:
		controller.NotFound(w, r)
	}
}
