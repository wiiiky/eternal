package user

type UpdateCoverRequest struct {
	Cover string `json:"cover" validate:"required"`
}
