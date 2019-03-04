pragma solidity >=0.4.18;

contract Food {

    struct PigInfo {
        address Creator;
        string id;
        string name;
        uint weight;
        uint date;
        string place;
    }

    struct FoodInfo {
        address Creator;
        string id;
        string name;
        uint volume;
        uint producedDate;
        uint packageDate;
        uint expireTime;
        string pigId;
        uint shopDate;
        uint score;
    }

    struct CheckInfo {
        address Creator;
        string foodId;
        uint checkDate;
        string checkRes;
        string checkDesc;
    }

    uint pig_num = 0;
    uint food_num = 0;
    uint chekcinfo_num = 0;

    mapping(uint => FoodInfo) FoodList; // key 为商品唯一编号
    mapping(uint => PigInfo) PigList;
    mapping(uint => CheckInfo) CheckInfoList;

    function utilCompareInternal(string memory a, string memory b) internal pure returns (bool) {
        if (bytes(a).length != bytes(b).length) {
            return false;
        }
        for (uint i = 0; i < bytes(a).length; i ++) {
            if (bytes(a)[i] != bytes(b)[i]) {
                return false;
            }
        }
        return true;
    }

    // PigInfo 相关操作 ========================================
    // 增加猪的信息
    function addPigInfo(address _creator, string memory _name, string memory _id, uint _weight, uint _date, string memory _place) public {
        PigInfo memory item = PigInfo(_creator, _id, _name, _weight, _date, _place);
        PigList[pig_num] = item;
        pig_num = pig_num + 1;
    }

    // 修改商品
    function updatePigInfoByIndex(string memory _id, string memory _name, uint _weight, uint _date, string memory _place) public {
        uint index = getPigIndexByID(_id);
        PigInfo memory item = PigList[index];
        item.name = _name;
        item.weight = _weight;
        item.date = _date;
        item.place = _place;
        PigList[index] = item;
    }

    function getPigIndexByID(string memory _id) public view returns (uint) {
        uint _index;
        for (uint index = 0; index < getPigNumber(); index++) {
            PigInfo memory item = PigList[index];
            bool isEqual = utilCompareInternal(_id, item.id);
            if (isEqual)
                _index = index;
        }
        return _index;
    }

    function getPigInfoByIndex(uint index) public view returns (address, string memory id, string memory name, uint weight, uint date, string memory place) {
        PigInfo memory item = PigList[index];
        return (item.Creator, item.id, item.name, item.weight, item.date, item.place);
    }

    function getPigInfoByID(string memory _id) public view returns (address, string memory id, string memory name, uint weight, uint date, string memory place) {
        uint index = getPigIndexByID(_id);
        return getPigInfoByIndex(index);
    }

    function getPigNumber() public view returns (uint) {
        return pig_num;
    }

    // FoodInfo 相关操作 ========================================
    // 添加商品信息
    function addFoodInfo(address _creator, string memory _id, string memory _name, uint _volume, uint _producedDate, uint _packageDate, uint _expireTime,
        string memory _pigId, uint _checkDate, string memory _checkRes, string memory _checkDesc, uint _shopDate, uint _score) public {
        FoodInfo memory item = FoodInfo(_creator, _id, _name, _volume, _producedDate, _packageDate, _expireTime, _pigId, _shopDate, _score);
        FoodList[food_num] = item;
        food_num = food_num + 1;
    }

    // 修改商品
    function updateFoodInfoByIndex(string memory _id, string memory _name, uint _volume, uint _producedDate, uint _packageDate, uint _expireTime, string memory _pigId,
        uint _checkDate, string memory _checkRes, string memory _checkDesc, uint _shopDate, uint _score) public {
        uint index = getFoodIndexByID(_id);
        FoodInfo memory item = FoodList[index];
        item.name = _name;
        item.volume = _volume;
        item.producedDate = _producedDate;
        item.packageDate = _packageDate;
        item.expireTime = _expireTime;
        item.pigId = _pigId;
        item.shopDate = _shopDate;
        item.score = _score;
        FoodList[index] = item;
    }

    function updateShopDate(string memory _id, uint _shopDate) public {
        uint index = getFoodIndexByID(_id);
        FoodInfo memory item = FoodList[index];
        item.shopDate = _shopDate;
        FoodList[index] = item;
    }

    function updateScore(string memory _id, uint _score) public {
        uint index = getFoodIndexByID(_id);
        FoodInfo memory item = FoodList[index];
        item.score = _score;
        FoodList[index] = item;
    }

    function getFoodIndexByID(string memory _id) public view returns (uint) {
        uint _index;
        for (uint index = 0; index < getFoodNumber(); index++) {
            FoodInfo memory item = FoodList[index];
            bool isEqual = utilCompareInternal(_id, item.id);
            if (isEqual)
                _index = index;
        }
        return _index;
    }

    function getFoodInfoByIndex(uint index) public view returns (address, string memory id, string memory name, uint volume,
        uint producedDate, uint packageDate, uint expireTime, string memory pigId, uint shopDate, uint score) {
        FoodInfo memory item = FoodList[index];
        return (item.Creator, item.id, item.name, item.volume, item.producedDate, item.packageDate, item.expireTime,
        item.pigId, item.shopDate, item.score);
    }

    function getFoodInfoByID(string memory _id) public view returns (address, string memory id, string memory name, uint volume,
        uint producedDate, uint packageDate, uint expireTime, string memory pigId, uint shopDate, uint score) {
        uint index = getFoodIndexByID(_id);
        return getFoodInfoByIndex(index);
    }

    function getFoodNumber() public view returns (uint) {
        return food_num;
    }

    function addCheckInfo(address _creator, string memory _foodId, uint _checkDate, string memory _checkRes, string memory _checkDesc) public {
        CheckInfo memory item = CheckInfo(_creator, _foodId, _checkDate, _checkRes, _checkDesc);
        CheckInfoList[chekcinfo_num] = item;
        chekcinfo_num++;
    }

    function getCheckInfoNumber() public view returns (uint) {
        return chekcinfo_num;
    }
}
