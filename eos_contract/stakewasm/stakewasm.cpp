#include "stakewasm.hpp"

namespace eosio {

void stakewasm::add( const uint64_t& tel, const asset& bonus, const asset& stake)
{
    require_auth( get_self() );
    check( std::to_string(tel).length() == 11, "invalid telephone" );
    check( bonus.is_valid(), "invalid bonus" );
    check( bonus.amount > 0, "must add positive bonus" );
    check( bonus.symbol == sym, "bonus symbol precision mismatch");
    check( stake.is_valid(), "invalid stake" );
    check( stake.amount == 0, "must stake zero asset" );
    check( stake.symbol == sym, "stake symbol precision mismatch");

    stakess ssinfo(get_self(), tel);
    auto it = ssinfo.find(tel);
    if(it != ssinfo.end()){
       ssinfo.modify( it, get_self(), [&]( auto& a ) {
          a.bonus += bonus;
       });
    }else{
       ssinfo.emplace(get_self(), [&](auto &a) {
          a.tel = tel;
          a.bonus = bonus;
          a.stake = stake;
       });
    }
}    

void stakewasm::sub( const uint64_t& tel, const asset& bonus)
{
    require_auth( get_self() );
    check( std::to_string(tel).length() == 11, "invalid telephone" );
    check( bonus.is_valid(), "invalid bonus" );
    check( bonus.amount > 0, "must sub positive unstake" );
    check( bonus.symbol == sym, "symbol precision mismatch");

    stakess ssinfo(get_self(), tel);
    auto it = ssinfo.find(tel);
    check(it != ssinfo.end(), "no such telephone in table");
    check(it->bonus.amount - it->stake.amount >= bonus.amount, "no much available balance to sub");
    ssinfo.modify( it, get_self(), [&]( auto& a ) {
        a.bonus -= bonus;
    });
}

void stakewasm::del( const uint64_t& tel)
{
    require_auth( get_self() );
    check( std::to_string(tel).length() == 11, "invalid telephone" );

    stakess ssinfo(get_self(), tel);
    auto it = ssinfo.find(tel);
    if(it != ssinfo.end()){
        ssinfo.erase(it);
    }
}

void stakewasm::stake( const uint64_t& tel, const asset& stake)
{
    require_auth( get_self() );
    check( std::to_string(tel).length() == 11, "invalid telephone" );
    check( stake.is_valid(), "invalid stake" );
    check( stake.amount > 0, "must transfer positive stake" );
    check( stake.symbol == sym, "symbol precision mismatch");

    stakess ssinfo(get_self(), tel);
    auto it = ssinfo.find(tel);
    check(it != ssinfo.end(), "no such telephone in table");
    check(it->stake.amount + stake.amount <= it->bonus.amount, "no much available balance to stake");
    ssinfo.modify( it, get_self(), [&]( auto& a ) {
        a.stake += stake;
    });
}

void stakewasm::unstake( const uint64_t& tel, const asset& unstake)
{
    require_auth( get_self() );
    check( std::to_string(tel).length() == 11, "invalid telephone" );
    check( unstake.is_valid(), "invalid unstake" );
    check( unstake.amount > 0, "must transfer positive unstake" );
    check( unstake.symbol == sym, "symbol precision mismatch");

    stakess ssinfo(get_self(), tel);
    auto it = ssinfo.find(tel);
    check(it != ssinfo.end(), "no such telephone in table");
    check(it->stake.amount >= unstake.amount, "no much available stake to unstake");
    ssinfo.modify( it, get_self(), [&]( auto& a ) {
        a.stake -= unstake;
    });
}

void stakewasm::show( const uint64_t& tel )
{
    check( std::to_string(tel).length() == 11, "invalid telephone");

    stakess ssinfo(get_self(), tel);
    auto it = ssinfo.find(tel);
    check(it != ssinfo.end(), "no sucn telephone in table");

    print("tel:", it->tel);
    print("bonus:", it->bonus.amount);
    print("stake:", it->stake.amount);
}
};
