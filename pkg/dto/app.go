package dto

type AddTaskRequest struct {
	URL          string `json:"url" form:"url" binding:"required"`
	Name         string `json:"name" form:"name" binding:"required"`
	TotalEp      int    `json:"total_ep" form:"total_ep" binding:"required"`
	CurrentEp    int    `json:"current_ep" form:"current_ep" binding:"required"`
	DownloadPath string `json:"download_path" form:"download_path" binding:"required"`
	Type         string `json:"type" form:"type" binding:"required"`
	Downloader   string `json:"downloader" form:"downloader" binding:"required"`
	Provider     string `json:"provider" form:"provider" binding:"required"`
}

type AddTaskResponse struct {
	ID int `json:"id"`
}

func (AddTaskRequest) GetMessages() bootstrap.ValidatorMessages {
	return bootstrap.ValidatorMessages{
		"url.required":           "URL不能为空",
		"name.required":          "名称不能为空",
		"total_ep.required":      "总集数不能为空",
		"current_ep.required":    "当前集数不能为空",
		"download_path.required": "下载路径不能为空",
		"type.required":          "类型不能为空",
		"downloader.required":    "下载器不能为空",
		"provider.required":      "提供商不能为空",
	}
}
