### 1. profile 수정
    - 패스워드 수정
    - 개인 정보 수정. 이메일, 닉네임
    - 회원 탈퇴

### 2. social 계정  
    - front에서는 회원가입이든 로그인이든 각 버튼 클릭은 동일 할 것이다.    
    - front에서는 provider, 해당하는 provider의 accessToken만을 backend로 전송한다.  
    - 구글, 페이스 북 developer 페이지에 회사 access ID를 등록, get 해야 한다.

### 3. 빌링
    - 충전
    - 아이템구매
        - 이력관리 중요
    - 충전 히스토리 /  구매 히스토리


### 4. API ... 
    - 서비스ID, 서비스 key를 통한 MD5로 인증 시행
        - smd5는 인풋 파라미더의 알파벳 순  + sid (즉 sid,skey는 항상 맨 끝)
            - ex. send data to backend 
                a_param1 = "a"
                b_param2 = "b"
                ...
                z_param3 = "z"
                ...
                sid = "sid"
                skey = "skey"
                smd5 = md5(a+b+...+z+...+sid+skey) // value of param with ordering alphabet

    - 사이버머니 충전 / 아이템 구매만 해당
        - externalUid 필요. 이 서비스를 이용하는 게임 혹은 서비스의 결제 추적용

###
    - 모든 ID 관련 변수는 대문자로 한다. 데이터베이스에는 ORM을 이용했기 때문에 i_d 형식으로 테이블이 만들어진다.
        - ex. UID => u_i_d

### package command
    govendor 
    govendor add +external