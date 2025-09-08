# Image Processing Pipeline in Go

## 📌 Overview

This project implements a **concurrent image processing pipeline** in Go. It downloads images from the web, applies a grayscale transformation, and stores both the original and processed versions locally.

The design leverages **goroutines, channels, WaitGroups, and mutexes** to build a safe, efficient, and scalable pipeline—demonstrating strong skills in **concurrency, synchronization, and system design** with Go.

---

## ✨ Features

* 🔗 **Concurrent pipeline** with separate stages: download → process → collect.
* 🎨 **Image transformation**: converts images to grayscale.
* 💾 **Storage**: saves original images to `input/` and processed images to `output/`.
* 🛡️ **Thread-safe collection** using a global map protected by a mutex.
* ⚡ **Efficient execution** with buffered and unbuffered channels.

---

## 🛠️ Tech Stack

* **Language**: Go (Golang)
* **Libraries**:

  * `image`, `image/color`, `image/png`, `image/jpeg` for image handling
  * `net/http` for downloading images
  * `sync` for concurrency control

---

## 📂 Project Structure

```
.
├── input/        # Stores downloaded original images
├── output/       # Stores processed grayscale images
├── main.go       # Entry point with pipeline implementation
```

---

## 🚀 Getting Started

### Prerequisites

* [Go](https://go.dev/dl/) 1.20+ installed on your system
* Internet connection (to fetch images)

### Run the project

```bash
# Clone the repository
git clone https://github.com/yourusername/go-image-pipeline.git
cd go-image-pipeline

# Run the program
go run main.go
```

### Output

* Original images → `input/`
* Grayscale images → `output/`
* Console logs show pipeline progress

---

## 💡 Why this project matters

This project demonstrates:

* Writing **clean, concurrent Go code**.
* Building a **producer-consumer pipeline** using channels.
* Applying **mutexes and WaitGroups** for synchronization.
* Handling **real-world tasks** like downloading, decoding, and processing images.

It serves as a **portfolio project** to highlight strong **backend and system-level programming skills**.

---

## 📜 License

This project is licensed under the MIT License – feel free to use and modify.

