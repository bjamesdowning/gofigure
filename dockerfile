FROM scratch
EXPOSE 8080
COPY gofigure /
COPY templates/ templates
COPY static/ static
CMD ["/gofigure"]
