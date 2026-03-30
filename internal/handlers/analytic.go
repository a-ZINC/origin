package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"origin.me/internal/store"
)

type AnalyticsHandler struct{ store *store.Store }

func NewAnalyticsHandler(st *store.Store) *AnalyticsHandler {
	return &AnalyticsHandler{store: st}
}

func (h *AnalyticsHandler) Overview(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.ListAllPosts(r.Context())
	if err != nil {
		fail(w, "db error", 500)
		return
	}

	type PostStat struct {
		Post          interface{} `json:"post"`
		Publications  interface{} `json:"publications"`
		SiteViews     int         `json:"siteViews"`
		TotalViews    int         `json:"totalViews"`
		TotalReact    int         `json:"totalReactions"`
		TotalComments int         `json:"totalComments"`
	}

	var stats []PostStat
	var grandViews, grandReact, grandComments int

	for _, p := range posts {
		if !p.Published {
			continue
		}
		pubs, _ := h.store.GetLatestAnalytics(r.Context(), p.ID)
		views, react, comments, _ := h.store.GetTotalAnalytics(r.Context(), p.ID)
		siteViews, _ := h.store.GetPostSiteViews(r.Context(), p.ID)

		total := views + siteViews
		grandViews += total
		grandReact += react
		grandComments += comments

		stats = append(stats, PostStat{
			Post:          p,
			Publications:  pubs,
			SiteViews:     siteViews,
			TotalViews:    total,
			TotalReact:    react,
			TotalComments: comments,
		})
	}

	siteChart, _ := h.store.GetSiteViewsLast14Days(r.Context())

	ok(w, map[string]interface{}{
		"posts":         stats,
		"siteChart":     siteChart,
		"grandViews":    grandViews,
		"grandReact":    grandReact,
		"grandComments": grandComments,
	})
}

func (h *AnalyticsHandler) PostDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := h.store.GetPostByID(r.Context(), id)
	if err != nil || post == nil {
		fail(w, "not found", 404)
		return
	}

	pubs, _ := h.store.GetLatestAnalytics(r.Context(), id)
	siteViews, _ := h.store.GetSiteAnalytics(r.Context(), id, 30)
	views, react, comments, _ := h.store.GetTotalAnalytics(r.Context(), id)
	siteTotal, _ := h.store.GetPostSiteViews(r.Context(), id)

	ok(w, map[string]interface{}{
		"post":          post,
		"publications":  pubs,
		"siteViews":     siteViews,
		"totalViews":    views + siteTotal,
		"totalReact":    react,
		"totalComments": comments,
	})
}

func (h *AnalyticsHandler) Track(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PostID string `json:"postId"`
	}
	if err := readJson(r, &body); err != nil || body.PostID == "" {
		fail(w, "postId required", 400)
		return
	}
	go h.store.TrackSiteView(r.Context(), body.PostID)
	w.WriteHeader(http.StatusNoContent)
}