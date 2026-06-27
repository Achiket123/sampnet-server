package research

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	domain "server/internal/domain/research"

	"github.com/gin-gonic/gin"
)

type mockUseCase struct {
	createFunc                 func(ctx context.Context, req domain.CreateRequest) (*domain.ResearchEntry, error)
	getByIDFunc                func(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error)
	listFunc                   func(ctx context.Context, req domain.ListRequest) ([]domain.ResearchEntry, int, error)
	updateFunc                 func(ctx context.Context, req domain.UpdateRequest, actorID uint, actorRole string, orgID uint) error
	deleteFunc                 func(ctx context.Context, id uint, actorID uint, actorRole string, orgID uint) error
	createFolderFunc           func(ctx context.Context, req domain.CreateFolderRequest, actorID uint, orgID uint) (*domain.ResearchFolder, error)
	getFolderByIDFunc          func(ctx context.Context, id uint, orgID uint) (*domain.ResearchFolder, error)
	getFoldersByResearchIDFunc func(ctx context.Context, researchID uint, parentID *uint, orgID uint) ([]domain.ResearchFolder, error)
	updateFolderFunc           func(ctx context.Context, req domain.UpdateFolderRequest, actorID uint, orgID uint) error
	deleteFolderFunc           func(ctx context.Context, id uint, actorID uint, orgID uint) error
	createDocFunc              func(ctx context.Context, req domain.CreateDocumentRequest, actorID uint, orgID uint) (*domain.ResearchDocument, error)
	getDocByIDFunc             func(ctx context.Context, id uint, orgID uint) (*domain.ResearchDocument, error)
	getDocsByFolderIDFunc      func(ctx context.Context, researchID uint, folderID *uint, orgID uint) ([]domain.ResearchDocument, error)
	updateDocFunc              func(ctx context.Context, req domain.UpdateDocumentRequest, actorID uint, orgID uint) error
	deleteDocFunc              func(ctx context.Context, id uint, actorID uint, orgID uint) error
	uploadFileFunc             func(ctx context.Context, req domain.CreateFileRequest, actorID uint, orgID uint) (*domain.ResearchFile, error)
	getFileByIDFunc            func(ctx context.Context, id uint, orgID uint) (*domain.ResearchFile, error)
	deleteFileFunc             func(ctx context.Context, id uint, actorID uint, orgID uint) error
	getFilesByDocIDFunc        func(ctx context.Context, docID uint, orgID uint) ([]domain.ResearchFile, error)
	addRefFunc                 func(ctx context.Context, req domain.CreateReferenceRequest, actorID uint, orgID uint) (*domain.ResearchReference, error)
	deleteRefFunc              func(ctx context.Context, id uint, actorID uint, orgID uint) error
	getRefsByDocIDFunc         func(ctx context.Context, docID uint, orgID uint) ([]domain.ResearchReference, error)
}

func (m *mockUseCase) Create(ctx context.Context, req domain.CreateRequest) (*domain.ResearchEntry, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, req)
	}
	return &domain.ResearchEntry{ID: 1, Title: req.Title}, nil
}

func (m *mockUseCase) GetByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id, orgID)
	}
	return &domain.ResearchEntry{ID: id, OrganisationID: orgID}, nil
}

func (m *mockUseCase) List(ctx context.Context, req domain.ListRequest) ([]domain.ResearchEntry, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, req)
	}
	return []domain.ResearchEntry{{ID: 1}}, 1, nil
}

func (m *mockUseCase) Update(ctx context.Context, req domain.UpdateRequest, actorID uint, actorRole string, orgID uint) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, req, actorID, actorRole, orgID)
	}
	return nil
}

func (m *mockUseCase) Delete(ctx context.Context, id uint, actorID uint, actorRole string, orgID uint) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id, actorID, actorRole, orgID)
	}
	return nil
}

func (m *mockUseCase) CreateFolder(ctx context.Context, req domain.CreateFolderRequest, actorID uint, orgID uint) (*domain.ResearchFolder, error) {
	if m.createFolderFunc != nil {
		return m.createFolderFunc(ctx, req, actorID, orgID)
	}
	return &domain.ResearchFolder{ID: 10, Name: req.Name}, nil
}

func (m *mockUseCase) GetFolderByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchFolder, error) {
	if m.getFolderByIDFunc != nil {
		return m.getFolderByIDFunc(ctx, id, orgID)
	}
	return &domain.ResearchFolder{ID: id}, nil
}

func (m *mockUseCase) GetFoldersByResearchID(ctx context.Context, researchID uint, parentID *uint, orgID uint) ([]domain.ResearchFolder, error) {
	if m.getFoldersByResearchIDFunc != nil {
		return m.getFoldersByResearchIDFunc(ctx, researchID, parentID, orgID)
	}
	return []domain.ResearchFolder{{ID: 10}}, nil
}

func (m *mockUseCase) UpdateFolder(ctx context.Context, req domain.UpdateFolderRequest, actorID uint, orgID uint) error {
	if m.updateFolderFunc != nil {
		return m.updateFolderFunc(ctx, req, actorID, orgID)
	}
	return nil
}

func (m *mockUseCase) DeleteFolder(ctx context.Context, id uint, actorID uint, orgID uint) error {
	if m.deleteFolderFunc != nil {
		return m.deleteFolderFunc(ctx, id, actorID, orgID)
	}
	return nil
}

func (m *mockUseCase) CreateDocument(ctx context.Context, req domain.CreateDocumentRequest, actorID uint, orgID uint) (*domain.ResearchDocument, error) {
	if m.createDocFunc != nil {
		return m.createDocFunc(ctx, req, actorID, orgID)
	}
	return &domain.ResearchDocument{ID: 100, Title: req.Title}, nil
}

func (m *mockUseCase) GetDocumentByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchDocument, error) {
	if m.getDocByIDFunc != nil {
		return m.getDocByIDFunc(ctx, id, orgID)
	}
	return &domain.ResearchDocument{ID: id}, nil
}

func (m *mockUseCase) GetDocumentsByFolderID(ctx context.Context, researchID uint, folderID *uint, orgID uint) ([]domain.ResearchDocument, error) {
	if m.getDocsByFolderIDFunc != nil {
		return m.getDocsByFolderIDFunc(ctx, researchID, folderID, orgID)
	}
	return []domain.ResearchDocument{{ID: 100}}, nil
}

func (m *mockUseCase) UpdateDocument(ctx context.Context, req domain.UpdateDocumentRequest, actorID uint, orgID uint) error {
	if m.updateDocFunc != nil {
		return m.updateDocFunc(ctx, req, actorID, orgID)
	}
	return nil
}

func (m *mockUseCase) DeleteDocument(ctx context.Context, id uint, actorID uint, orgID uint) error {
	if m.deleteDocFunc != nil {
		return m.deleteDocFunc(ctx, id, actorID, orgID)
	}
	return nil
}

func (m *mockUseCase) UploadFile(ctx context.Context, req domain.CreateFileRequest, actorID uint, orgID uint) (*domain.ResearchFile, error) {
	if m.uploadFileFunc != nil {
		return m.uploadFileFunc(ctx, req, actorID, orgID)
	}
	return &domain.ResearchFile{ID: 1000, FileName: req.FileName}, nil
}

func (m *mockUseCase) GetFileByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchFile, error) {
	if m.getFileByIDFunc != nil {
		return m.getFileByIDFunc(ctx, id, orgID)
	}
	return &domain.ResearchFile{ID: id, StoragePath: "test.txt"}, nil
}

func (m *mockUseCase) DeleteFile(ctx context.Context, id uint, actorID uint, orgID uint) error {
	if m.deleteFileFunc != nil {
		return m.deleteFileFunc(ctx, id, actorID, orgID)
	}
	return nil
}

func (m *mockUseCase) GetFilesByDocumentID(ctx context.Context, docID uint, orgID uint) ([]domain.ResearchFile, error) {
	if m.getFilesByDocIDFunc != nil {
		return m.getFilesByDocIDFunc(ctx, docID, orgID)
	}
	return []domain.ResearchFile{{ID: 1000}}, nil
}

func (m *mockUseCase) AddReference(ctx context.Context, req domain.CreateReferenceRequest, actorID uint, orgID uint) (*domain.ResearchReference, error) {
	if m.addRefFunc != nil {
		return m.addRefFunc(ctx, req, actorID, orgID)
	}
	return &domain.ResearchReference{ID: 500, Title: req.Title}, nil
}

func (m *mockUseCase) DeleteReference(ctx context.Context, id uint, actorID uint, orgID uint) error {
	if m.deleteRefFunc != nil {
		return m.deleteRefFunc(ctx, id, actorID, orgID)
	}
	return nil
}

func (m *mockUseCase) GetReferencesByDocumentID(ctx context.Context, docID uint, orgID uint) ([]domain.ResearchReference, error) {
	if m.getRefsByDocIDFunc != nil {
		return m.getRefsByDocIDFunc(ctx, docID, orgID)
	}
	return []domain.ResearchReference{{ID: 500}}, nil
}

func setupTestContext(method, url string, body string, auth bool) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	c.Request = req
	if auth {
		c.Set("userID", uint(1))
		c.Set("organisationID", uint(10))
		c.Set("role", "user")
	}
	return w, c
}

func TestHandler_Create(t *testing.T) {
	uc := &mockUseCase{}
	h := NewHandler(uc)

	t.Run("unauthorized", func(t *testing.T) {
		w, c := setupTestContext(http.MethodPost, "/api/v1/research", `{"title":"Res"}`, false)
		h.Create(c)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		w, c := setupTestContext(http.MethodPost, "/api/v1/research", `invalid json`, true)
		h.Create(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		ucErr := &mockUseCase{
			createFunc: func(ctx context.Context, req domain.CreateRequest) (*domain.ResearchEntry, error) {
				return nil, errors.New("db error")
			},
		}
		hErr := NewHandler(ucErr)
		w, c := setupTestContext(http.MethodPost, "/api/v1/research", `{"title":"Res"}`, true)
		hErr.Create(c)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodPost, "/api/v1/research", `{"title":"Res","status":"draft"}`, true)
		h.Create(c)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})
}

func TestHandler_GetByID(t *testing.T) {
	uc := &mockUseCase{}
	h := NewHandler(uc)

	t.Run("invalid id", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/abc", "", true)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		h.GetByID(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		ucNotFound := &mockUseCase{
			getByIDFunc: func(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
				return nil, domain.ErrNotFound
			},
		}
		hNotFound := NewHandler(ucNotFound)
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/1", "", true)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		hNotFound.GetByID(c)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/1", "", true)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		h.GetByID(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestHandler_List(t *testing.T) {
	uc := &mockUseCase{}
	h := NewHandler(uc)

	t.Run("success with filters", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research?project_id=2&team_id=3&limit=10&offset=0", "", true)
		h.List(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestHandler_UpdateAndDelete(t *testing.T) {
	uc := &mockUseCase{}
	h := NewHandler(uc)

	t.Run("Update forbidden", func(t *testing.T) {
		ucForbidden := &mockUseCase{
			updateFunc: func(ctx context.Context, req domain.UpdateRequest, actorID uint, actorRole string, orgID uint) error {
				return domain.ErrUnauthorized
			},
		}
		hForbidden := NewHandler(ucForbidden)
		w, c := setupTestContext(http.MethodPut, "/api/v1/research/1", `{"title":"Updated"}`, true)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		hForbidden.Update(c)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("Update success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodPut, "/api/v1/research/1", `{"title":"Updated"}`, true)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		h.Update(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Delete success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodDelete, "/api/v1/research/1", "", true)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		h.Delete(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestHandler_FoldersAndDocuments(t *testing.T) {
	uc := &mockUseCase{}
	h := NewHandler(uc)

	t.Run("CreateFolder success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodPost, "/api/v1/research/folders", `{"research_id":1,"name":"Docs"}`, true)
		h.CreateFolder(c)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})

	t.Run("ListFolders missing research_id", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/folders", "", true)
		h.ListFolders(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("ListFolders success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/folders?research_id=1&parent_id=10", "", true)
		h.ListFolders(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("CreateDocument success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodPost, "/api/v1/research/documents", `{"research_id":1,"title":"Doc"}`, true)
		h.CreateDocument(c)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})

	t.Run("ListDocuments success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/documents?research_id=1&folder_id=10", "", true)
		h.ListDocuments(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestHandler_FilesAndReferences(t *testing.T) {
	uc := &mockUseCase{}
	h := NewHandler(uc)

	t.Run("ListFiles missing document_id", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/files", "", true)
		h.ListFiles(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("ListFiles success", func(t *testing.T) {
		w, c := setupTestContext(http.MethodGet, "/api/v1/research/files?document_id=100", "", true)
		h.ListFiles(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("AddReference and ListReferences", func(t *testing.T) {
		w, c := setupTestContext(http.MethodPost, "/api/v1/research/references", `{"research_id":1,"title":"Ref"}`, true)
		h.AddReference(c)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}

		w2, c2 := setupTestContext(http.MethodGet, "/api/v1/research/references?document_id=100", "", true)
		h.ListReferences(c2)
		if w2.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w2.Code)
		}
	})

	t.Run("UploadFile success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("research_id", "1")
		_ = writer.WriteField("document_id", "100")
		_ = writer.WriteField("folder_id", "10")
		part, _ := writer.CreateFormFile("file", "test.txt")
		_, _ = part.Write([]byte("hello world"))
		_ = writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/research/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		c.Request = req
		c.Set("userID", uint(1))
		c.Set("organisationID", uint(10))
		c.Set("role", "user")

		h.UploadFile(c)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Cleanup uploads directory created during test
		_ = os.RemoveAll("uploads")
	})

	t.Run("DownloadFile success", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "sample.txt")
		_ = os.WriteFile(tmpFile, []byte("test content"), 0644)

		ucDownload := &mockUseCase{
			getFileByIDFunc: func(ctx context.Context, id uint, orgID uint) (*domain.ResearchFile, error) {
				return &domain.ResearchFile{ID: id, StoragePath: tmpFile}, nil
			},
		}
		hDownload := NewHandler(ucDownload)

		w, c := setupTestContext(http.MethodGet, "/api/v1/research/files/1000/download", "", true)
		c.Params = gin.Params{{Key: "id", Value: "1000"}}
		hDownload.DownloadFile(c)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}
