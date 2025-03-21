

sudo apt install build-essential


go_version=$(go version 2>/dev/null)
required_go_version="1.18"
target_go_version="1.24.1"

go_dl_url="https://go.dev/dl/go"$target_go_version".linux-amd64.tar.gz"
echo "Kiểm tra Go version..."
if [ -z "$go_version" ]; then
    echo "Go chưa được cài đặt. Tiến hành cài đặt Go..."
    wget $go_dl_url
    sudo tar -C /usr/local -xzf go"$target_go_version".linux-amd64.tar.gz
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
    echo "Go đã được cài đặt thành công!"
else
    installed_go_version=$(echo $go_version | awk '{print $3}' | cut -d '.' -f1,2)
    if [[ "$installed_go_version" < "$required_go_version" ]]; then
        echo "Go phiên bản $installed_go_version không đủ, cài đặt lại..."
        wget $go_dl_url
        sudo tar -C /usr/local -xzf go"$target_go_version".linux-amd64.tar.gz
        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
        source ~/.bashrc
        echo "Go đã được cài đặt lại thành công!"
    else
        echo "Go đã được cài đặt và đủ phiên bản!"
    fi
fi

# Kiểm tra phiên bản Rust
rust_version=$(rustc --version 2>/dev/null)
required_rust_version="1.81.0"

echo "Kiểm tra Rust version..."
if [ -z "$rust_version" ]; then
    echo "Rust chưa được cài đặt. Tiến hành cài đặt Rust..."
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
    source $HOME/.cargo/env
    echo "Rust đã được cài đặt thành công!"
else
    installed_rust_version=$(echo $rust_version | awk '{print $2}')
    if [[ "$installed_rust_version" < "$required_rust_version" ]]; then
        echo "Rust phiên bản $installed_rust_version không đủ, cài đặt lại..."
        rustup self uninstall -y
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
        source $HOME/.cargo/env
        
        echo "Rust đã được cài đặt lại thành công!"
    else
        echo "Rust đã được cài đặt và đủ phiên bản!"
    fi
fi
source ~/.bashrc


# Kiểm tra xem lệnh rzup có sẵn không
if command -v rzup &>/dev/null; then
    echo "rzup đã được cài đặt và sẵn sàng sử dụng!"
    rzup --version
else
    echo "rzup chưa được cài đặt. Tiến hành cài đặt..."

    # Cài đặt rzup từ risczero.com
    curl -L https://risczero.com/install | bash

    # Thêm thư mục chứa rzup vào PATH
    echo "Đang thêm /root/.risc0/bin vào PATH..."
    echo "export PATH=\$PATH:/root/.risc0/bin" >> ~/.bashrc

    # Áp dụng thay đổi PATH
    source ~/.bashrc

    # Kiểm tra lại sau khi cài đặt
    if command -v rzup &>/dev/null; then
        echo "Cài đặt rzup thành công!"
        rzup --version
    else
        echo "Cài đặt rzup không thành công. Vui lòng kiểm tra lại!"
    fi
fi

rzup install


echo "Kiểm tra và cài đặt hoàn tất!"
