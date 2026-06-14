package handler

import (
	"context"

	"github.com/noedaka/go-config-parser/cmd/grpc_server/internal/proto"
	"github.com/noedaka/go-config-parser/internal/parser"
	"github.com/noedaka/go-config-parser/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	proto.UnimplementedConfigRecsServiceServer

	rules  []service.Rule
	parser parser.Parser
}

func NewHandler(rules []service.Rule, parser parser.Parser) *handler {
	return &handler{rules: rules, parser: parser}
}

func (h *handler) ConfigRecommendationsByFileHandler(
	ctx context.Context, r *proto.ConfigFileRequest,
) (*proto.RecsResponse, error) {
	data := r.GetData()

	cfg, err := h.parser.ParseConfig(data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot parse data")
	}

	var issues []service.Issue
	for _, rule := range h.rules {
		issues = append(issues, rule.Check(cfg)...)
	}

	recs := service.FormatIssues(issues)
	response := proto.RecsResponse_builder{
		Recs: &recs,
	}

	return response.Build(), nil
}

func (h *handler) ConfigRecommendationsByStringHandler(
	ctx context.Context, r *proto.ConfigStringRequest,
) (*proto.RecsResponse, error) {
	data := r.GetData()

	cfg, err := h.parser.ParseConfig([]byte(data))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot parse data")
	}

	var issues []service.Issue
	for _, rule := range h.rules {
		issues = append(issues, rule.Check(cfg)...)
	}

	recs := service.FormatIssues(issues)
	response := proto.RecsResponse_builder{
		Recs: &recs,
	}

	return response.Build(), nil
}
