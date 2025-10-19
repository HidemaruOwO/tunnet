## **Bridge Tunnel Protocol 要約**

### **目的**

TCP/UDP通信を**ChaCha20-Poly1305**で暗号化し、**Curve25519 (ECDH)** で鍵交換する軽量高速トンネル。ゲーム・リモート制御など低遅延用途向け。

---

### **通信構成**

```
Client App ←→ bridge tunnel ←→ [暗号化TCP/UDP] ←→ bridge server ←→ Target
```

---

### **暗号スイート（固定）**

| 要素   | アルゴリズム      | 詳細                        |
| ------ | ----------------- | --------------------------- |
| 鍵交換 | Curve25519 (ECDH) | 毎接続時に一時鍵生成        |
| 鍵導出 | HKDF-SHA256       | 送受信方向で独立            |
| 暗号化 | ChaCha20-Poly1305 | AEAD (認証付き暗号)         |
| ノンス | 12B               | `0x00000000 \| counter/seq` |

---

### **ハンドシェイク**

```
1. 公開鍵交換: [pub_len(1B)] [pub]
2. 共有秘密: shared = X25519(priv, peer_pub)
3. 鍵導出:
   sessionID = SHA256(shared || client_pub || server_pub)
   c2sKey = HKDF(shared, salt=sessionID, info="bridge c2s")
   s2cKey = HKDF(shared, salt=sessionID, info="bridge s2c")
4. AAD: "dir:c2s," || sessionID (方向別)
```

---

### **TCPフレーム**

```
[len(4B)] [nonce(12B)] [ciphertext] [auth_tag(16B)]
```

- **len**: 平文長 (暗号化しない)
- **nonce**: 8Bカウンタ + 4Bゼロ (方向別管理)
- タグ検証失敗時は即切断

---

### **UDPフレーム**

```
[seq(8B)] [nonce(12B)] [ciphertext] [auth_tag(16B)]
```

- **seq**: 64bitカウンタ (リプレイ防止)
- **nonce**: `0x00000000 | seq`
- **AAD**: `"dir:c2s," || sessionID || ",seq:" || seq`
- 重複/過去seqは破棄、再送制御なし

---

### **セキュリティ機構**

| 保護対象       | 実装                      |
| -------------- | ------------------------- |
| 機密性         | ChaCha20暗号化            |
| 完全性         | Poly1305タグ (改ざん検知) |
| リプレイ防止   | ノンス/seq管理            |
| セッション分離 | sessionID                 |
| 方向混在防止   | AAD(dir)                  |

---

## **CLI例**

以下、詳細な説明を残したまま構造化して要約しました：

## Bridge コマンドリファレンス

### サーバーモード (`bridge server`)

#### 基本起動

```bash
bridge server  # デフォルト: 0.0.0.0:7070、認証なし
bridge server -c config.toml  # 設定ファイル指定 (-c: --config)
bridge server -addr 201.234.213.56 -port 17171  # アドレス・ポート指定
```

- `-addr` (--bind-address): バインドアドレス
- `-port` (--bind-port): バインドポート

#### 認証方式

```bash
bridge server --token <token>  # トークン認証
bridge server --password <hash or password> --algorithm argon2  # パスワード認証
bridge server --public-key id25519.pub  # 公開鍵認証
```

- パスワードアルゴリズム例: `bcrypt`, `scrypt`, `argon2`, `plain`

---

### クライアントモード (`bridge tunnel`)

#### TCP トンネル

```bash
bridge tunnel tcp 3000 -addr 201.234.213.56  # ローカル3000番ポートを公開
bridge tunnel tcp 25565 -addr mc.server.com -fast  # 低遅延モード（非暗号化）
bridge tunnel tcp 8080 -addr your.server.com -port 17171  # カスタムポート指定
```

- `-addr` (--server-address): 接続先サーバーアドレス **【必須】**
- `-port` (--server-port): サーバーポート指定
- `-fast` (--no-encryption): 非暗号化トンネル（VPN間接続など低遅延用途向け）

#### 認証付き接続

```bash
bridge tunnel tcp 25565 -addr 10.1.0.202 --token <token>  # トークン認証
bridge tunnel tcp 25565 -addr 100.23.45.21 --password <password>  # パスワード認証
bridge tunnel tcp 25565 -addr 192.168.1.2 --private-key id25519  # 秘密鍵認証
```

#### UDP トンネル

```bash
bridge tunnel udp 53 -addr dns.server.com  # ローカル53番ポートを公開
bridge tunnel udp 19132 -addr bedrock.mc.server.com  # Minecraft Bedrock版など
```

#### 実装見送り機能

```bash
bridge tunnel both 8080 -addr app.server.com  # TCP+UDP同時トンネル（後で実装予定）
```

---

### シークレット管理 (`bridge secret`)

```bash
bridge secret gen-keypair  # 公開鍵・秘密鍵ペア生成
bridge secret hash-password <password> --algorithm argon2  # パスワードハッシュ生成
```

---

**主な特徴:**

- TCP/UDP両対応のトンネリングツール
- 複数認証方式（トークン/パスワード/公開鍵）
- 低遅延モード（非暗号化オプション）
- 設定ファイルまたはCLI引数での柔軟な設定

---

### **実装 (Go)**

| 機能 | パッケージ                             |
| ---- | -------------------------------------- |
| ECDH | `crypto/ecdh`                          |
| HKDF | `crypto/hkdf`                          |
| AEAD | `golang.org/x/crypto/chacha20poly1305` |
| Hash | `crypto/sha256`                        |

---

### **運用指針**

1. **セッション毎に鍵交換** (鍵再利用禁止)
2. **ノンス/seq再利用禁止** (上限前に再接続)
3. **タグ検証失敗時は即破棄**
4. UDPは順序保証なし
5. sessionIDはログ時にトランケート推奨

---

### **将来拡張案**

- PSK認証 (MITM防止)
- 0-RTT再接続
- Multiplex (複数トンネル多重化)
- ACK付きUDP
- 固定長パディング (長さ秘匿)

---

**参考**: RFC 7539 (ChaCha20/Poly1305), RFC 7748 (X25519), RFC 5869 (HKDF)
