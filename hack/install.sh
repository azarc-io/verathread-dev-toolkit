#!/usr/bin/env bash

APP_NAME="vdt"
REPO_URL="https://vth-artifacts.s3-eu-west-1.amazonaws.com/vdt"

: ${VTH_INSTALL_DIR:="/usr/local/bin"}
: ${USE_SUDO:="true"}

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="x86_64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(uname|tr '[:upper:]' '[:upper:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
}

# runs the given command as root (detects if we are root already)
runAsRoot() {
  local CMD="$*"

  if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
    CMD="sudo $CMD"
  fi

  $CMD
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  local supported="Darwin-x86_64\nLinux-x86_64\nDarwin-arm64\nLinux-arm64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuilt binary for ${OS}-${ARCH}."
    echo "To build from source, go to $REPO_URL"
    exit 1
  fi

  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# checkTagProvided checks whether TAG has provided as an environment variable so we can skip checkLatestVersion.
checkTagProvided() {
  [[ ! -z "$TAG" ]]
}

# checkLatestVersion grabs the latest version string from the releases
checkLatestVersion() {
  local latest_release_url="$REPO_URL/latest.txt"
  if type "curl" > /dev/null; then
    echo "Checking ${latest_release_url} for latest release"
    TAG=$(curl -Ls $latest_release_url | grep -oE "[^/]+$" )
  elif type "wget" > /dev/null; then
    TAG=$(wget -O - $latest_release_url > /dev/null 2>&1 | grep -oE "[^/]+$")
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  VTH_DIST="vth-cli-$TAG-$OS-$ARCH.tar.gz"
  DOWNLOAD_URL="$REPO_URL/v$TAG/$VTH_DIST"
  VTH_TMP_ROOT="$(mktemp -dt vdt-binary-XXXXXX)"
  VTH_TMP_FILE="$VTH_TMP_ROOT/$VTH_DIST"
  echo "Downloading ${DOWNLOAD_URL}"
  if type "curl" > /dev/null; then
    curl -SsL "$DOWNLOAD_URL" -o "$VTH_TMP_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$VTH_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  echo "Preparing to install $APP_NAME into ${VTH_INSTALL_DIR}"
  tar -xf "$VTH_TMP_FILE"
  rm -f "$VTH_INSTALL_DIR/$APP_NAME"
  runAsRoot chmod 755 "$APP_NAME"
  runAsRoot cp "$APP_NAME" "$VTH_INSTALL_DIR/$APP_NAME"
  echo "$APP_NAME installed into $VTH_INSTALL_DIR/$APP_NAME"
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    if [[ -n "$INPUT_ARGUMENTS" ]]; then
      echo "Failed to install $APP_NAME with the arguments provided: $INPUT_ARGUMENTS"
      help
    else
      echo "Failed to install $APP_NAME"
    fi
    echo -e "\tFor support, go to $REPO_URL."
  fi
  cleanup
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  if ! command -v $APP_NAME &> /dev/null; then
    echo "$APP_NAME not found. Is $VTH_INSTALL_DIR on your "'$PATH?'
    exit 1
  fi
  echo "Run '$APP_NAME --help' to see what you can do with it."
}

# checks which version of the cli is installed and
# if it needs to be changed.
checkInstalledVersion() {
  if [[ -f "${VTH_INSTALL_DIR}/${APP_NAME}" ]]; then
    local version=$(${APP_NAME} version)
    if [[ "$version" == "$TAG" ]]; then
      echo "${APP_NAME} ${version} is already ${DESIRED_VERSION:-latest}"
      return 0
    else
      echo "${APP_NAME} ${TAG} is available. Changing from version ${version}."
      return 1
    fi
  else
    return 1
  fi
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  if ! command -v $APP_NAME &> /dev/null; then
    echo "$APP_NAME not found. Is $VTH_INSTALL_DIR on your "'$PATH?'
    exit 1
  fi
  echo "Run '$APP_NAME --help' to see what you can do with it."
}

# cleanup temporary files
cleanup() {
  if [[ -d "${VTH_TMP_ROOT:-}" ]]; then
    rm -rf "$VTH_TMP_ROOT"
  fi
}

#Stop execution on any error
trap "fail_trap" EXIT
set -e

initArch
initOS
verifySupported
checkTagProvided || checkLatestVersion
if ! checkInstalledVersion; then
  downloadFile
  installFile
fi
testVersion
cleanup
