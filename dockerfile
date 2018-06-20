FROM scratch
EXPOSE 8080
COPY gofigure /
COPY templates/ templates
COPY templates/css templates/css
CMD ["/gofigure"]
