FROM scratch
EXPOSE 8080
COPY gofigure /
COPY templates/ templates
CMD ["/gofigure"]
