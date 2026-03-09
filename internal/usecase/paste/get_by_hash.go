package paste

import (
	"context"
	"errors"
	"fmt"

	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
	views "github.com/StewardMcCormick/Paste_Bin/internal/util/views_worker"
)

func (uc *UseCase) GetByHash(ctx context.Context, request dto.GetPasteRequest, hash string) (*dto.PasteResponse, error) {
	log := appctx.GetLogger(ctx)

	paste, err := uc.repo.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, errs.PasteNotFound) {
			return nil, err
		}
		log.Error(fmt.Sprintf("%v - get paste error", err))
		return nil, fmt.Errorf("%w - get paste error", errs.InternalError)
	}

	userId, err := appctx.GetUserId(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("%v - get user_id from ctx error", err))
		return nil, fmt.Errorf("%w - get user_id from ctx error", errs.InternalError)
	}

	if paste.Privacy == domain.PrivatePolicy && userId != paste.UserId {
		log.Debug(fmt.Sprintf("get paste Forbidden(Private): from - %d, to paste with user_id - %d", userId, paste.UserId))
		return nil, errs.Forbidden
	} else if paste.Privacy == domain.ProtectedPolicy &&
		!uc.security.CompareHashAndPassword(paste.PasswordHash, request.Password) {
		log.Debug(fmt.Sprintf("get paste Forbidden(Protected): from - %d, to paste with user_id - %d", userId, paste.UserId))
		return nil, errs.Unauthorized
	}

	paste.Views += 1
	uc.viewWorker.SendEvent(ctx, views.ViewEvent{PasteId: paste.Id})

	return paste.ToResponse(), nil
}
