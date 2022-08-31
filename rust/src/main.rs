use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::TcpListener;
use tokio::sync::broadcast;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let listener = TcpListener::bind("127.0.0.1:1337").await?;

    let (tx, _rx) = broadcast::channel::<Vec<u8>>(100);

    loop {
        let (mut conn, _) = listener.accept().await?;

        let tx = tx.clone();
        let mut rx = tx.subscribe();
        tokio::spawn(async move {
            let (reader, mut writer) = conn.split();
            let mut reader = tokio::io::BufReader::new(reader);
            let mut buf = [0u8; 4096];

            loop {
                tokio::select! {
                    res = reader.read(&mut buf) => {
                        if let Ok(n) = res {
                            if n == 0 {
                                break;
                            }
                            tx.send(buf[..n].to_vec()).unwrap();
                        }
                    },
                    res = rx.recv() => {
                        if let Ok(msg) = res {
                            writer.write_all(&msg).await.unwrap();
                        }
                    }
                }
            }
        });
    }
}
