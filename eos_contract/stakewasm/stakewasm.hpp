#pragma once

#include <eosio/asset.hpp>
#include <eosio/symbol.hpp>
#include <eosio/eosio.hpp>

#include <string>

namespace eosiosystem {
   class system_contract;
}

namespace eosio {

   using std::string;

      class [[eosio::contract("stakewasm")]] stakewasm : public contract {
      public:
         using contract::contract;

         [[eosio::action]]
         void add( const uint64_t& tel, const asset& bonus, const asset& stake);

         [[eosio::action]]
         void sub( const uint64_t& tel, const asset& bonus);

         [[eosio::action]]
         void del( const uint64_t& tel);

         [[eosio::action]]
         void stake( const uint64_t& tel, const asset& stake);

         [[eosio::action]]
         void unstake( const uint64_t& tel, const asset& stake);

         [[eosio::action]]
         void show( const uint64_t& tel );

         using add_action = eosio::action_wrapper<"add"_n, &stakewasm::add>;
         using sub_action = eosio::action_wrapper<"sub"_n, &stakewasm::sub>;
         using del_action = eosio::action_wrapper<"del"_n, &stakewasm::del>;
         using stake_action = eosio::action_wrapper<"stake"_n, &stakewasm::stake>;
         using unstake_action = eosio::action_wrapper<"unstake"_n, &stakewasm::unstake>;
         using show_action = eosio::action_wrapper<"show"_n, &stakewasm::show>;

      private:
         struct [[eosio::table]] stake_status {
            uint64_t    tel;
            asset       bonus;
            asset       stake;

            uint64_t primary_key() const { return tel; }
         };  

         typedef eosio::multi_index<"stakess"_n, stake_status > stakess;  

         symbol sym = symbol(symbol_code("MSN"), 4);

      };
   /** @}*/ // end of @defgroup 
} /// namespace eosio          