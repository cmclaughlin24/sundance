package handlers

import "net/http"

func (h *Handlers) GetCanonicalTags(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetCanonicalTag(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) CreateCanonicalTag(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeleteCanonicalTag(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetCanonicalTagVersions(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) CreateCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) PublishCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeprecateCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) RetireCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}
