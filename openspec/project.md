# 椤圭洰涓婁笅鏂?

## Purpose
CharonOMS 鏄竴涓暀鑲插煿璁満鏋勭殑瀛︾敓绠＄悊 SaaS 绯荤粺锛屼粠鍘?Python Flask 椤圭洰 (ZhixinStudentSaaS) 杩佺Щ閲嶅啓涓?Go 鐗堟湰銆傜郴缁熶富瑕佺敤浜庯細
- 瀛︾敓淇℃伅绠＄悊锛堝鐢熸。妗堛€佸勾绾с€佽仈绯绘柟寮忥級
- 鏁欑粌淇℃伅绠＄悊锛堟暀缁冩。妗堛€佹巿璇惧绉戙€佽仈绯绘柟寮忥級
- 瀛︾敓-鏁欑粌鍏宠仈绠＄悊锛堝瀵瑰鍏崇郴锛?
- 璁㈠崟绠＄悊锛堣绋嬭鍗曘€佹敹娆俱€侀€€娆撅級
- 鍩轰簬瑙掕壊鐨勬潈闄愭帶鍒讹紙RBAC锛夌郴缁?
- 鑿滃崟鍜屾潈闄愬姩鎬佺鐞?

**椤圭洰鐩爣**锛?
1. 浣跨敤 Go 鎻愬崌绯荤粺鎬ц兘鍜屽苟鍙戝鐞嗚兘鍔?
2. 閲囩敤 DDD锛堥鍩熼┍鍔ㄨ璁★級鏋舵瀯鎻愰珮浠ｇ爜鍙淮鎶ゆ€?
3. 淇濇寔涓庡師 Python 鐗堟湰鐨?API 鍏煎鎬э紝纭繚鍓嶇鏃犵紳杩佺Щ
4. 寤虹珛娓呮櫚鐨勫垎灞傛灦鏋勫拰渚濊禆鍏崇郴

## 鎶€鏈爤

### 鍚庣
- Go 1.21+
- Gin v1.9.1 (Web 妗嗘灦)
- GORM v1.25.5 (ORM)
- MySQL 5.7+ (鏁版嵁搴擄紝utf8mb4 瀛楃闆?
- JWT (golang-jwt/jwt v5.2.0, 璁よ瘉)
- Zap v1.26.0 (鏃ュ織)
- Viper v1.18.2 (閰嶇疆绠＄悊)
- bcrypt (瀵嗙爜鍔犲瘑)

### 鍓嶇
- Vue.js 3 (鍘熼」鐩墠绔紝鏈慨鏀?
- Axios (HTTP 瀹㈡埛绔?
- 鍘熺敓 HTML/CSS/JavaScript

### 鏁版嵁搴撹〃
- `useraccount` - 鐢ㄦ埛璐﹀彿
- `role` - 瑙掕壊
- `menu` - 鑿滃崟
- `permissions` - 鏉冮檺
- `role_permissions` - 瑙掕壊-鏉冮檺鍏宠仈
- `sex` - 鎬у埆瀛楀吀
- `grade` - 骞寸骇瀛楀吀
- `subject` - 瀛︾瀛楀吀
- `student` - 瀛︾敓淇℃伅
- `coach` - 鏁欑粌淇℃伅
- `student_coach` - 瀛︾敓-鏁欑粌鍏宠仈锛堝瀵瑰锛?

## 椤圭洰绾﹀畾

### 浠ｇ爜椋庢牸
- **Go 瑙勮寖**: 涓ユ牸閬靛惊 Go 瀹樻柟浠ｇ爜瑙勮寖锛屼娇鐢?`gofmt` 鍜?`goimports` 鏍煎紡鍖?
- **鍛藉悕绾﹀畾**:
  - 鏂囦欢鍚嶏細灏忓啓涓嬪垝绾?(snake_case)锛屽 `auth_handler.go`
  - 鍖呭悕锛氬皬鍐欏崟璇嶏紝濡?`package auth`
  - 瀵煎嚭鍑芥暟/绫诲瀷锛氬ぇ椹煎嘲 (PascalCase)锛屽 `NewAuthService`
  - 绉佹湁鍑芥暟/瀛楁锛氬皬椹煎嘲 (camelCase)锛屽 `hashPassword`
  - 甯搁噺锛氬ぇ鍐欎笅鍒掔嚎锛屽 `MAX_RETRY_COUNT`
- **娉ㄩ噴**: 鎵€鏈夊鍑虹殑鍑芥暟銆佺被鍨嬨€佸父閲忓繀椤绘坊鍔犳敞閲婏紝娉ㄩ噴浠ヨ娉ㄩ噴瀵硅薄鍚嶇О寮€澶?
- **閿欒澶勭悊**: 閿欒蹇呴』鏄惧紡妫€鏌ワ紝浣跨敤 `pkg/errors` 鍖呭畾涔変笟鍔￠敊璇?

### 鏋舵瀯妯″紡

**DDD 鍒嗗眰鏋舵瀯**锛堜粠澶栧埌鍐咃級锛?

1. **鎺ュ彛灞?(Interfaces)** - `internal/interfaces/`
   - HTTP Handler: 澶勭悊 HTTP 璇锋眰/鍝嶅簲
   - DTO: 鏁版嵁浼犺緭瀵硅薄
   - Middleware: 璁よ瘉銆佹棩蹇椼€丆ORS
   - Router: 璺敱閰嶇疆

2. **搴旂敤灞?(Application)** - `internal/application/`
   - Service: 缂栨帓棰嗗煙瀵硅薄锛屽疄鐜颁笟鍔＄敤渚?
   - Assembler: DTO 鈫?Entity 杞崲

3. **棰嗗煙灞?(Domain)** - `internal/domain/`
   - Entity: 棰嗗煙瀹炰綋锛堝寘鍚笟鍔￠€昏緫锛?
   - Repository Interface: 浠撳偍鎺ュ彛瀹氫箟
   - Domain Service: 棰嗗煙鏈嶅姟

4. **鍩虹璁炬柦灞?(Infrastructure)** - `internal/infrastructure/`
   - Persistence: 浠撳偍瀹炵幇锛圡ySQL锛?
   - Config: 閰嶇疆鍔犺浇
   - Logger: 鏃ュ織瀹炵幇

**渚濊禆瑙勫垯**锛?
- 渚濊禆鏂瑰悜锛欼nterfaces 鈫?Application 鈫?Domain 鈫?Infrastructure
- Domain 灞備笉渚濊禆浠讳綍鍏朵粬灞傦紙渚濊禆鍊掔疆锛?
- Infrastructure 灞傚疄鐜?Domain 灞傚畾涔夌殑鎺ュ彛

**妯″潡缁勭粐**锛氭寜涓氬姟妯″潡鍒掑垎锛坄auth`, `rbac`, `basic`, `student`, `coach`, `order`锛夛紝姣忎釜妯″潡鍖呭惈瀹屾暣鐨勫洓灞傜粨鏋?

### 娴嬭瘯绛栫暐
- **鍗曞厓娴嬭瘯**: 浣跨敤 `testing` 鍖咃紝瑕嗙洊鏍稿績涓氬姟閫昏緫
- **闆嗘垚娴嬭瘯**: 娴嬭瘯 Repository 涓庢暟鎹簱浜や簰
- **娴嬭瘯鏂囦欢**: `*_test.go`
- **瑕嗙洊鐜囩洰鏍?*: 棰嗗煙灞?> 80%锛屽簲鐢ㄥ眰 > 70%
- **Mock**: 浣跨敤鎺ュ彛杩涜渚濊禆娉ㄥ叆

### Git 宸ヤ綔娴?

**鍒嗘敮绛栫暐**锛?
- `main` - 涓诲垎鏀紙鐢熶骇鐜锛?
- `develop` - 寮€鍙戝垎鏀?
- `feature/*` - 鍔熻兘鍒嗘敮
- `fix/*` - 淇鍒嗘敮

**鎻愪氦瑙勮寖** (Conventional Commits)锛?
```
<type>(<scope>): <subject>
```

**Type**锛歚feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

**绀轰緥**锛?
```
feat(auth): implement JWT authentication

- Add JWT token generation and validation
- Implement login and logout handlers
- Add JWT middleware for protected routes
```

## 棰嗗煙涓婁笅鏂?

### 涓氬姟棰嗗煙锛氭暀鑲插煿璁涓?

**鏍稿績姒傚康**锛?
- **瀛︾敓 (Student)**: 鎺ュ彈鍩硅鐨勫鍛橈紝鍖呭惈鍩烘湰淇℃伅銆佸勾绾с€佽仈绯绘柟寮?
- **鏁欑粌 (Coach)**: 鎺堣鏁欏笀锛屽寘鍚熀鏈俊鎭€佹巿璇惧绉戙€佽仈绯绘柟寮?
- **瀛︾ (Subject)**: 鎺堣绉戠洰锛堣鏂囥€佹暟瀛︺€佽嫳璇瓑锛?
- **骞寸骇 (Grade)**: 瀛︾敓鎵€鍦ㄥ勾绾э紙涓€骞寸骇~楂樹笁锛?
- **璁㈠崟 (Order)**: 璇剧▼璐拱璁㈠崟
- **瀛︾敓-鏁欑粌鍏崇郴**: 澶氬澶氬叧绯伙紝涓€涓鐢熷彲浠ユ湁澶氫釜鏁欑粌锛屼竴涓暀缁冨彲浠ユ暀澶氫釜瀛︾敓

**鏉冮檺绠＄悊 (RBAC)**锛?
- **鐢ㄦ埛璐﹀彿 (UserAccount)**: 绯荤粺鐧诲綍鐢ㄦ埛
- **瑙掕壊 (Role)**: 濡?瓒呯骇绠＄悊鍛?銆?鏅€氱鐞嗗憳"銆?鏁欏姟"绛?
- **鑿滃崟 (Menu)**: 绯荤粺鍔熻兘妯″潡鐨勬爲褰㈢粨鏋?
- **鏉冮檺 (Permission)**: 瀵硅彍鍗曠殑鎿嶄綔鏉冮檺
- **瑙掕壊-鏉冮檺缁戝畾**: 閫氳繃 `role_permissions` 琛ㄥ叧鑱?

**涓氬姟瑙勫垯**锛?
- 鐢ㄦ埛蹇呴』鍏宠仈瑙掕壊鎵嶈兘鐧诲綍绯荤粺
- 鍒犻櫎瀛︾敓鍓嶉渶妫€鏌ユ槸鍚︽湁鍏宠仈璁㈠崟锛堟湁璁㈠崟鍒欑姝㈠垹闄わ級
- 鍒犻櫎鏁欑粌鏃惰嚜鍔ㄨВ闄や笌瀛︾敓鐨勫叧鑱斿叧绯?
- 鑿滃崟閲囩敤鏍戝舰缁撴瀯锛岄《灞傝彍鍗曠殑 `parent_id` 涓?NULL
- 鐘舵€佸瓧娈电粺涓€浣跨敤 `status`: 0=鍚敤, 1=绂佺敤

**鏁版嵁鐘舵€?*锛?
- **鍚敤/绂佺敤**: `status` 瀛楁锛?-鍚敤锛?-绂佺敤锛?
- **杞垹闄?*: 閮ㄥ垎琛ㄤ娇鐢?`deleted_at` 瀛楁锛圙ORM 杞垹闄わ級
- **鏃堕棿鎴?*: `create_time`, `update_time` 鑷姩缁存姢

## 閲嶈绾︽潫

### 鎶€鏈害鏉?
1. **UTF-8 缂栫爜**:
   - 鎵€鏈夋簮浠ｇ爜鏂囦欢蹇呴』浣跨敤 UTF-8 鏃?BOM 缂栫爜
   - 鏁版嵁搴撲娇鐢?utf8mb4 瀛楃闆?
   - PowerShell 鑴氭湰淇敼鏂囦欢鏃跺繀椤绘樉寮忔寚瀹?UTF-8 缂栫爜

2. **API 鍏煎鎬?*:
   - 蹇呴』涓庡師 Python 鐗堟湰鐨?API 鍝嶅簲鏍煎紡淇濇寔涓€鑷?
   - 鍝嶅簲鏍煎紡锛歚{"code": 0, "message": "success", "data": {...}}`
   - 鍓嶇鏈慨鏀癸紝鍚庣蹇呴』瀹屽叏鍏煎鍓嶇璋冪敤

3. **鏁版嵁搴撶害鏉?*:
   - 澶栭敭绾︽潫锛歚menu.parent_id` 寮曠敤 `menu.id`锛堥《灞備负 NULL锛屼笉鑳界敤 0锛?
   - 瀛楁鍛藉悕锛氫笌鍘熸暟鎹簱淇濇寔涓€鑷达紙濡?`comment` 鑰岄潪 `description`锛?

4. **鎬ц兘绾︽潫**:
   - 浣跨敤杩炴帴姹犵鐞嗘暟鎹簱杩炴帴
   - JWT token 楠岃瘉蹇呴』楂樻晥锛堥伩鍏嶉绻佹煡搴擄級

### 涓氬姟绾︽潫
- 瓒呯骇绠＄悊鍛?(`is_super_admin = 1`) 鎷ユ湁鎵€鏈夋潈闄愶紝璺宠繃鏉冮檺妫€鏌?
- 瑙掕壊鍒犻櫎鍓嶉渶妫€鏌ユ槸鍚︽湁鐢ㄦ埛鍏宠仈
- 瀛︾敓鍒犻櫎鍓嶉渶妫€鏌ユ槸鍚︽湁璁㈠崟鍏宠仈

### 瀹夊叏绾︽潫
- 瀵嗙爜浣跨敤 bcrypt 鍔犲瘑锛屾垚鏈洜瀛?10
- JWT Secret 蹇呴』閰嶇疆鍦ㄧ幆澧冨彉閲忔垨閰嶇疆鏂囦欢涓?
- 鏁忔劅鎺ュ彛蹇呴』閫氳繃 JWT 涓棿浠朵繚鎶?
- CORS 閰嶇疆闄愬埗鍏佽鐨勬潵婧?

## 澶栭儴渚濊禆

### 鏁版嵁搴?
- **MySQL 8.0+**
  - Host: 閰嶇疆鍦?`config/config.yaml`
  - 瀛楃闆? utf8mb4_unicode_ci
  - 杩炴帴姹? 鏈€澶ц繛鎺ユ暟 100锛屾渶澶х┖闂茶繛鎺?10

### 閰嶇疆鏂囦欢
- **config/config.yaml**: 涓婚厤缃枃浠?
  - 鏁版嵁搴撹繛鎺ヤ俊鎭?
  - JWT secret 鍜岃繃鏈熸椂闂?
  - 鏈嶅姟鍣ㄧ鍙ｅ拰杩愯妯″紡
  - CORS 閰嶇疆

### 鍓嶇璧勬簮
- 鍓嶇闈欐€佹枃浠朵綅浜?`frontend/` 鐩綍
- 閫氳繃 Gin 鐨?`Static` 鍜?`StaticFile` 鎻愪緵鏈嶅姟
- 璺緞鏄犲皠锛歚/` 鈫?`frontend/index.html`, `/frontend/*` 鈫?`frontend/*`

### 寮€鍙戞敞鎰忎簨椤?

**甯歌闄烽槺**锛?
1. PowerShell 淇敼鏂囦欢鏃跺繀椤讳娇鐢?`[System.Text.UTF8Encoding]::new($false)` 鎸囧畾 UTF-8 鏃?BOM
2. GORM 浣跨敤 `gorm` 鏍囩鏄庣‘鎸囧畾鍒楀悕锛岄伩鍏嶅懡鍚嶈鍒欒嚜鍔ㄨ浆鎹?
3. JWT 涓棿浠跺繀椤诲湪闇€瑕佽璇佺殑璺敱缁勪箣鍓嶆敞鍐?
4. HTTP handler 涓崟鑾锋墍鏈夐敊璇紝浣跨敤缁熶竴鐨勫搷搴旀牸寮?
5. 娑夊強澶氳〃鎿嶄綔鏃朵娇鐢?GORM 浜嬪姟

**杩佺Щ妫€鏌ユ竻鍗?*锛堜粠鍘?Python 鐗堟湰杩佺Щ鍔熻兘锛夛細
- [ ] API 璺緞鍜屾柟娉曚竴鑷?
- [ ] 璇锋眰/鍝嶅簲 JSON 鏍煎紡涓€鑷?
- [ ] 鏁版嵁搴撹〃缁撴瀯鍜屽瓧娈靛悕涓€鑷?
- [ ] 涓氬姟閫昏緫鍜岄獙璇佽鍒欎竴鑷?
- [ ] 閿欒鍝嶅簲鏍煎紡鍜岀姸鎬佺爜涓€鑷?
- [ ] 娴嬭瘯鍓嶇璋冪敤姝ｅ父宸ヤ綔