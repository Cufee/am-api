# Use the official image as a parent image.
FROM golang

# Set the working directory.
WORKDIR /app

# Copy the file from your host to your current location.
COPY ./build .

# Run the command inside your image filesystem.
RUN chmod +x app

# Add metadata to the image to describe which port the container is listening on at runtime.
EXPOSE 4000

# Run the specified command within the container.
CMD [ "./app" ]